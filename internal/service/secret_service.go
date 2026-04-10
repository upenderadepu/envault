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

type SecretService struct {
	secretMetaRepo *repository.SecretMetadataRepository
	envRepo        *repository.EnvironmentRepository
	projectRepo    *repository.ProjectRepository
	auditRepo      *repository.AuditLogRepository
	vaultSvc       *vault.VaultService
}

func NewSecretService(
	secretMetaRepo *repository.SecretMetadataRepository,
	envRepo *repository.EnvironmentRepository,
	projectRepo *repository.ProjectRepository,
	auditRepo *repository.AuditLogRepository,
	vaultSvc *vault.VaultService,
) *SecretService {
	return &SecretService{
		secretMetaRepo: secretMetaRepo,
		envRepo:        envRepo,
		projectRepo:    projectRepo,
		auditRepo:      auditRepo,
		vaultSvc:       vaultSvc,
	}
}

// ListKeys returns secret metadata for a project/environment. No values returned.
func (s *SecretService) ListKeys(projectSlug, envName string) ([]models.SecretMetadata, error) {
	project, err := s.projectRepo.FindBySlug(projectSlug)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	env, err := s.envRepo.FindByProjectAndName(project.ID, envName)
	if err != nil {
		return nil, fmt.Errorf("environment not found: %w", err)
	}

	return s.secretMetaRepo.FindByProjectAndEnv(project.ID, env.ID)
}

// GetSecret reads a single secret VALUE from Vault. This is the ONLY method
// that returns actual secret values.
func (s *SecretService) GetSecret(ctx context.Context, projectSlug, envName, keyName string, userID uuid.UUID) (string, *models.SecretMetadata, error) {
	project, err := s.projectRepo.FindBySlug(projectSlug)
	if err != nil {
		return "", nil, fmt.Errorf("project not found: %w", err)
	}

	env, err := s.envRepo.FindByProjectAndName(project.ID, envName)
	if err != nil {
		return "", nil, fmt.Errorf("environment not found: %w", err)
	}

	meta, err := s.secretMetaRepo.FindByKey(project.ID, env.ID, keyName)
	if err != nil {
		return "", nil, fmt.Errorf("secret key not found: %w", err)
	}

	// Read from Vault
	data, err := s.vaultSvc.ReadSecret(ctx, project.VaultMountPath, envName)
	if err != nil {
		return "", nil, fmt.Errorf("vault read failed: %w", err)
	}

	val, ok := data[keyName]
	if !ok {
		return "", nil, fmt.Errorf("key %s not found in vault", keyName)
	}

	// Audit log (key name only, NEVER the value)
	auditMeta, _ := json.Marshal(map[string]string{"key": keyName, "environment": envName})
	s.auditRepo.Create(&models.AuditLog{
		ProjectID:    project.ID,
		UserID:       &userID,
		Action:       models.ActionSecretRead,
		ResourcePath: fmt.Sprintf("projects/%s/environments/%s/secrets/%s", projectSlug, envName, keyName),
		Metadata:     auditMeta,
	})

	return fmt.Sprintf("%v", val), meta, nil
}

// SetSecret writes a single key-value pair to Vault (read-merge-write) and upserts metadata.
func (s *SecretService) SetSecret(ctx context.Context, projectSlug, envName, keyName, value string, userID uuid.UUID) (*models.SecretMetadata, error) {
	project, err := s.projectRepo.FindBySlug(projectSlug)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	env, err := s.envRepo.FindByProjectAndName(project.ID, envName)
	if err != nil {
		return nil, fmt.Errorf("environment not found: %w", err)
	}

	// Write to Vault (read-merge-write internally)
	if err := s.vaultSvc.WriteSecret(ctx, project.VaultMountPath, envName, map[string]interface{}{keyName: value}); err != nil {
		return nil, fmt.Errorf("vault write failed: %w", err)
	}

	// Upsert metadata
	meta := &models.SecretMetadata{
		ProjectID:      project.ID,
		EnvironmentID:  env.ID,
		KeyName:        keyName,
		VaultPath:      fmt.Sprintf("%s/%s", project.VaultMountPath, envName),
		CreatedByID:    userID,
		VaultVersion:   1,
		LastModifiedAt: time.Now(),
	}

	// Try to find existing to increment version
	existing, _ := s.secretMetaRepo.FindByKey(project.ID, env.ID, keyName)
	if existing != nil {
		meta.VaultVersion = existing.VaultVersion + 1
	}

	if err := s.secretMetaRepo.Upsert(meta); err != nil {
		return nil, fmt.Errorf("upsert metadata: %w", err)
	}

	// Audit log (key name only, NEVER the value)
	auditMeta, _ := json.Marshal(map[string]string{"key": keyName, "environment": envName})
	s.auditRepo.Create(&models.AuditLog{
		ProjectID:    project.ID,
		UserID:       &userID,
		Action:       models.ActionSecretWrite,
		ResourcePath: fmt.Sprintf("projects/%s/environments/%s/secrets/%s", projectSlug, envName, keyName),
		Metadata:     auditMeta,
	})

	// Reload
	meta, _ = s.secretMetaRepo.FindByKey(project.ID, env.ID, keyName)
	return meta, nil
}

// DeleteSecret removes a single key from Vault and deletes its metadata.
func (s *SecretService) DeleteSecret(ctx context.Context, projectSlug, envName, keyName string, userID uuid.UUID) error {
	project, err := s.projectRepo.FindBySlug(projectSlug)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	env, err := s.envRepo.FindByProjectAndName(project.ID, envName)
	if err != nil {
		return fmt.Errorf("environment not found: %w", err)
	}

	if err := s.vaultSvc.DeleteSecretKey(ctx, project.VaultMountPath, envName, keyName); err != nil {
		return fmt.Errorf("vault delete failed: %w", err)
	}

	if err := s.secretMetaRepo.DeleteByKey(project.ID, env.ID, keyName); err != nil {
		return fmt.Errorf("delete metadata: %w", err)
	}

	auditMeta, _ := json.Marshal(map[string]string{"key": keyName, "environment": envName})
	s.auditRepo.Create(&models.AuditLog{
		ProjectID:    project.ID,
		UserID:       &userID,
		Action:       models.ActionSecretDelete,
		ResourcePath: fmt.Sprintf("projects/%s/environments/%s/secrets/%s", projectSlug, envName, keyName),
		Metadata:     auditMeta,
	})

	return nil
}

// BulkSetSecrets writes multiple key-value pairs in a single Vault write.
func (s *SecretService) BulkSetSecrets(ctx context.Context, projectSlug, envName string, secrets map[string]string, userID uuid.UUID) ([]models.SecretMetadata, error) {
	project, err := s.projectRepo.FindBySlug(projectSlug)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	env, err := s.envRepo.FindByProjectAndName(project.ID, envName)
	if err != nil {
		return nil, fmt.Errorf("environment not found: %w", err)
	}

	// Convert to map[string]interface{} for Vault
	data := make(map[string]interface{}, len(secrets))
	for k, v := range secrets {
		data[k] = v
	}

	if err := s.vaultSvc.WriteSecret(ctx, project.VaultMountPath, envName, data); err != nil {
		return nil, fmt.Errorf("vault bulk write failed: %w", err)
	}

	// Upsert metadata for each key
	now := time.Now()
	for keyName := range secrets {
		existing, _ := s.secretMetaRepo.FindByKey(project.ID, env.ID, keyName)
		version := 1
		if existing != nil {
			version = existing.VaultVersion + 1
		}

		meta := &models.SecretMetadata{
			ProjectID:      project.ID,
			EnvironmentID:  env.ID,
			KeyName:        keyName,
			VaultPath:      fmt.Sprintf("%s/%s", project.VaultMountPath, envName),
			CreatedByID:    userID,
			VaultVersion:   version,
			LastModifiedAt: now,
		}
		s.secretMetaRepo.Upsert(meta)
	}

	// Audit log
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	auditMeta, _ := json.Marshal(map[string]interface{}{"keys": keys, "environment": envName, "count": len(secrets)})
	s.auditRepo.Create(&models.AuditLog{
		ProjectID:    project.ID,
		UserID:       &userID,
		Action:       models.ActionSecretWrite,
		ResourcePath: fmt.Sprintf("projects/%s/environments/%s/secrets", projectSlug, envName),
		Metadata:     auditMeta,
	})

	return s.secretMetaRepo.FindByProjectAndEnv(project.ID, env.ID)
}
