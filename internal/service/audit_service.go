package service

import (
	"fmt"

	"github.com/bhartiyaanshul/envault/internal/models"
	"github.com/bhartiyaanshul/envault/internal/repository"
)

type AuditService struct {
	auditRepo   *repository.AuditLogRepository
	projectRepo *repository.ProjectRepository
}

func NewAuditService(
	auditRepo *repository.AuditLogRepository,
	projectRepo *repository.ProjectRepository,
) *AuditService {
	return &AuditService{auditRepo: auditRepo, projectRepo: projectRepo}
}

// ListAuditLogs returns paginated audit logs, optionally filtered by action.
func (s *AuditService) ListAuditLogs(projectSlug, action string, limit, offset int) ([]models.AuditLog, int64, error) {
	project, err := s.projectRepo.FindBySlug(projectSlug)
	if err != nil {
		return nil, 0, fmt.Errorf("project not found: %w", err)
	}

	if action != "" {
		return s.auditRepo.FindByProjectAndAction(project.ID, action, limit, offset)
	}
	return s.auditRepo.FindByProjectID(project.ID, limit, offset)
}
