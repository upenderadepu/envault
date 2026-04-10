package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bhartiyaanshul/envault/internal/models"
	"github.com/bhartiyaanshul/envault/internal/repository"
	"github.com/bhartiyaanshul/envault/internal/vault"
	"github.com/google/uuid"
)

type MemberService struct {
	memberRepo  *repository.TeamMemberRepository
	userRepo    *repository.UserRepository
	projectRepo *repository.ProjectRepository
	auditRepo   *repository.AuditLogRepository
	vaultSvc    *vault.VaultService
}

func NewMemberService(
	memberRepo *repository.TeamMemberRepository,
	userRepo *repository.UserRepository,
	projectRepo *repository.ProjectRepository,
	auditRepo *repository.AuditLogRepository,
	vaultSvc *vault.VaultService,
) *MemberService {
	return &MemberService{
		memberRepo:  memberRepo,
		userRepo:    userRepo,
		projectRepo: projectRepo,
		auditRepo:   auditRepo,
		vaultSvc:    vaultSvc,
	}
}

func (s *MemberService) ListMembers(projectSlug string) ([]models.TeamMember, error) {
	project, err := s.projectRepo.FindBySlug(projectSlug)
	if err != nil {
		return nil, err
	}
	return s.memberRepo.FindActiveByProjectID(project.ID)
}

// AddMember invites a user to a project with a specific role.
// Returns the team member and a one-time Vault token.
func (s *MemberService) AddMember(ctx context.Context, projectSlug, email, role string, inviterID uuid.UUID) (*models.TeamMember, string, error) {
	project, err := s.projectRepo.FindBySlug(projectSlug)
	if err != nil {
		return nil, "", fmt.Errorf("project not found: %w", err)
	}

	user, err := s.userRepo.FindOrCreate("pending-"+email, email)
	if err != nil {
		return nil, "", fmt.Errorf("find/create user: %w", err)
	}

	// Determine environment access based on role
	envNames := environmentsForRole(role)

	// Build and write Vault policy
	policyName := fmt.Sprintf("%s-%s-%s", project.Slug, role, user.ID.String()[:8])
	policyHCL := s.vaultSvc.BuildPolicies(project.VaultMountPath, role, envNames)
	if err := s.vaultSvc.WritePolicy(ctx, policyName, policyHCL); err != nil {
		return nil, "", fmt.Errorf("write policy: %w", err)
	}

	// Token TTLs: CI gets shorter TTL
	ttl, maxTTL := tokenTTLForRole(role)

	token, accessor, err := s.vaultSvc.CreateUserToken(ctx, []string{policyName}, ttl, maxTTL)
	if err != nil {
		return nil, "", fmt.Errorf("create token: %w", err)
	}

	now := time.Now()
	member := &models.TeamMember{
		ProjectID:          project.ID,
		UserID:             user.ID,
		Role:               role,
		VaultPolicyName:    policyName,
		VaultTokenAccessor: accessor,
		IsActive:           true,
		JoinedAt:           &now,
	}

	if err := s.memberRepo.Create(member); err != nil {
		return nil, "", fmt.Errorf("create member: %w", err)
	}

	auditMeta, _ := json.Marshal(map[string]string{"email": email, "role": role})
	s.auditRepo.Create(&models.AuditLog{
		ProjectID:    project.ID,
		UserID:       &inviterID,
		Action:       models.ActionMemberInvite,
		ResourcePath: fmt.Sprintf("projects/%s/members/%s", projectSlug, user.ID),
		Metadata:     auditMeta,
	})

	member.User = user
	return member, token, nil
}

// RemoveMember deactivates a member and revokes their Vault token.
func (s *MemberService) RemoveMember(ctx context.Context, projectSlug string, memberID, removerID uuid.UUID) error {
	project, err := s.projectRepo.FindBySlug(projectSlug)
	if err != nil {
		return err
	}

	member, err := s.memberRepo.FindByID(memberID)
	if err != nil {
		return fmt.Errorf("member not found: %w", err)
	}

	// Revoke Vault token
	if err := s.vaultSvc.RevokeTokenByAccessor(ctx, member.VaultTokenAccessor); err != nil {
		return fmt.Errorf("revoke token: %w", err)
	}

	member.IsActive = false
	member.VaultTokenAccessor = ""
	if err := s.memberRepo.Update(member); err != nil {
		return fmt.Errorf("deactivate member: %w", err)
	}

	auditMeta, _ := json.Marshal(map[string]string{"member_id": memberID.String()})
	s.auditRepo.Create(&models.AuditLog{
		ProjectID:    project.ID,
		UserID:       &removerID,
		Action:       models.ActionMemberRemove,
		ResourcePath: fmt.Sprintf("projects/%s/members/%s", projectSlug, memberID),
		Metadata:     auditMeta,
	})

	return nil
}

// RotateCredentials revokes the old token and creates a new one with the same policies.
func (s *MemberService) RotateCredentials(ctx context.Context, projectSlug string, userID uuid.UUID) (string, error) {
	project, err := s.projectRepo.FindBySlug(projectSlug)
	if err != nil {
		return "", err
	}

	member, err := s.memberRepo.FindByProjectAndUser(project.ID, userID)
	if err != nil {
		return "", fmt.Errorf("member not found: %w", err)
	}

	// Revoke old token
	s.vaultSvc.RevokeTokenByAccessor(ctx, member.VaultTokenAccessor)

	// Create new token with same policy
	ttl, maxTTL := tokenTTLForRole(member.Role)
	token, accessor, err := s.vaultSvc.CreateUserToken(ctx, []string{member.VaultPolicyName}, ttl, maxTTL)
	if err != nil {
		return "", fmt.Errorf("create new token: %w", err)
	}

	member.VaultTokenAccessor = accessor
	if err := s.memberRepo.Update(member); err != nil {
		return "", fmt.Errorf("update accessor: %w", err)
	}

	s.auditRepo.Create(&models.AuditLog{
		ProjectID:    project.ID,
		UserID:       &userID,
		Action:       models.ActionCredentialsRotate,
		ResourcePath: fmt.Sprintf("projects/%s/rotate", projectSlug),
	})

	return token, nil
}

func environmentsForRole(role string) []string {
	switch role {
	case "admin":
		return []string{"development", "staging", "production"}
	case "developer":
		return []string{"development", "staging"}
	case "ci":
		return []string{"staging", "production"}
	default:
		return []string{"development"}
	}
}

func tokenTTLForRole(role string) (ttl, maxTTL time.Duration) {
	if role == "ci" {
		return 1 * time.Hour, 8 * time.Hour
	}
	return 8 * time.Hour, 24 * time.Hour
}
