package repository

import (
	"github.com/bhartiyaanshul/envault/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditLogRepository is append-only. There are no Update or Delete methods.
type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

// FindByProjectID returns paginated audit logs ordered by most recent first.
func (r *AuditLogRepository) FindByProjectID(projectID uuid.UUID, limit, offset int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	base := r.db.Model(&models.AuditLog{}).Where("project_id = ?", projectID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := base.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	return logs, total, err
}

// FindByProjectAndAction filters audit logs by action type.
func (r *AuditLogRepository) FindByProjectAndAction(projectID uuid.UUID, action string, limit, offset int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	base := r.db.Model(&models.AuditLog{}).Where("project_id = ? AND action = ?", projectID, action)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := base.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	return logs, total, err
}
