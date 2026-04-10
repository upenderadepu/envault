package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bhartiyaanshul/envault/internal/models"
	"github.com/bhartiyaanshul/envault/internal/repository"
	"github.com/bhartiyaanshul/envault/internal/vault"
	"github.com/google/uuid"
)

type ProjectService struct {
	projectRepo    *repository.ProjectRepository
	envRepo        *repository.EnvironmentRepository
	memberRepo     *repository.TeamMemberRepository
	auditRepo      *repository.AuditLogRepository
	secretMetaRepo *repository.SecretMetadataRepository
	vaultSvc       *vault.VaultService
	mountPrefix    string
}

func NewProjectService(
	projectRepo *repository.ProjectRepository,
	envRepo *repository.EnvironmentRepository,
	memberRepo *repository.TeamMemberRepository,
	auditRepo *repository.AuditLogRepository,
	secretMetaRepo *repository.SecretMetadataRepository,
	vaultSvc *vault.VaultService,
	mountPrefix string,
) *ProjectService {
	return &ProjectService{
		projectRepo:    projectRepo,
		envRepo:        envRepo,
		memberRepo:     memberRepo,
		auditRepo:      auditRepo,
		secretMetaRepo: secretMetaRepo,
		vaultSvc:       vaultSvc,
		mountPrefix:    mountPrefix,
	}
}

var slugRegex = regexp.MustCompile(`[^a-z0-9]+`)

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = slugRegex.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > 50 {
		slug = slug[:50]
	}
	return slug
}

// CreateProject orchestrates project creation:
// DB record -> environments -> Vault mount -> policy -> owner token -> team member -> audit.
// Returns the project and a one-time Vault token for the owner.
func (s *ProjectService) CreateProject(ctx context.Context, name string, ownerID uuid.UUID) (*models.Project, string, error) {
	slug := generateSlug(name)
	mountPath := fmt.Sprintf("%s-%s", s.mountPrefix, slug)

	project := &models.Project{
		Name:           name,
		Slug:           slug,
		VaultMountPath: mountPath,
		OwnerID:        ownerID,
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, "", fmt.Errorf("create project: %w", err)
	}

	// Create default environments
	envs := []models.Environment{
		{ProjectID: project.ID, Name: "development", IsProduction: false},
		{ProjectID: project.ID, Name: "staging", IsProduction: false},
		{ProjectID: project.ID, Name: "production", IsProduction: true},
	}
	if err := s.envRepo.CreateBatch(envs); err != nil {
		return nil, "", fmt.Errorf("create environments: %w", err)
	}

	// Enable KV-v2 mount in Vault
	if err := s.vaultSvc.EnableKV2(ctx, mountPath); err != nil {
		return nil, "", fmt.Errorf("enable kv2: %w", err)
	}

	// Build and write admin policy
	policyName := fmt.Sprintf("%s-admin", slug)
	policyHCL := s.vaultSvc.BuildPolicies(mountPath, "admin", []string{"development", "staging", "production"})
	if err := s.vaultSvc.WritePolicy(ctx, policyName, policyHCL); err != nil {
		return nil, "", fmt.Errorf("write policy: %w", err)
	}

	// Create Vault token for owner: 8h TTL, 24h max
	token, accessor, err := s.vaultSvc.CreateUserToken(ctx, []string{policyName}, 8*time.Hour, 24*time.Hour)
	if err != nil {
		return nil, "", fmt.Errorf("create owner token: %w", err)
	}

	// Record owner as admin team member (store accessor, NEVER the token)
	now := time.Now()
	member := &models.TeamMember{
		ProjectID:          project.ID,
		UserID:             ownerID,
		Role:               "admin",
		VaultPolicyName:    policyName,
		VaultTokenAccessor: accessor,
		IsActive:           true,
		JoinedAt:           &now,
	}
	if err := s.memberRepo.Create(member); err != nil {
		return nil, "", fmt.Errorf("create owner membership: %w", err)
	}

	// Audit log
	meta, _ := json.Marshal(map[string]string{"slug": slug})
	s.auditRepo.Create(&models.AuditLog{
		ProjectID:    project.ID,
		UserID:       &ownerID,
		Action:       models.ActionProjectCreate,
		ResourcePath: fmt.Sprintf("projects/%s", slug),
		Metadata:     meta,
	})

	// Reload with associations
	project, _ = s.projectRepo.FindBySlug(slug)
	return project, token, nil
}

func (s *ProjectService) GetProject(slug string) (*models.Project, error) {
	return s.projectRepo.FindBySlug(slug)
}

func (s *ProjectService) ListProjects(userID uuid.UUID) ([]models.Project, error) {
	return s.projectRepo.ListForUser(userID)
}

// DeleteProject revokes all Vault credentials and soft-deletes the project.
func (s *ProjectService) DeleteProject(ctx context.Context, slug string, userID uuid.UUID) error {
	project, err := s.projectRepo.FindBySlug(slug)
	if err != nil {
		return err
	}

	if project.OwnerID != userID {
		return fmt.Errorf("only the project owner can delete")
	}

	// Revoke all Vault tokens
	accessors, _ := s.memberRepo.GetAccessorsByProjectID(project.ID)
	if len(accessors) > 0 {
		s.vaultSvc.RevokeAllProjectCredentials(ctx, accessors)
	}

	if err := s.projectRepo.Delete(project.ID); err != nil {
		return err
	}

	meta, _ := json.Marshal(map[string]string{"slug": slug})
	s.auditRepo.Create(&models.AuditLog{
		ProjectID:    project.ID,
		UserID:       &userID,
		Action:       models.ActionProjectDelete,
		ResourcePath: fmt.Sprintf("projects/%s", slug),
		Metadata:     meta,
	})

	return nil
}
