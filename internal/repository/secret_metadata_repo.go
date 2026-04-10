package repository

import (
	"time"

	"github.com/bhartiyaanshul/envault/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SecretMetadataRepository struct {
	db *gorm.DB
}

func NewSecretMetadataRepository(db *gorm.DB) *SecretMetadataRepository {
	return &SecretMetadataRepository{db: db}
}

func (r *SecretMetadataRepository) Create(meta *models.SecretMetadata) error {
	return r.db.Create(meta).Error
}

// Upsert inserts or updates secret metadata on conflict (project_id, environment_id, key_name).
func (r *SecretMetadataRepository) Upsert(meta *models.SecretMetadata) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "project_id"},
			{Name: "environment_id"},
			{Name: "key_name"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"vault_version",
			"last_modified_at",
			"created_by_id",
		}),
	}).Create(meta).Error
}

func (r *SecretMetadataRepository) FindByProjectAndEnv(projectID, envID uuid.UUID) ([]models.SecretMetadata, error) {
	var metas []models.SecretMetadata
	err := r.db.Where("project_id = ? AND environment_id = ?", projectID, envID).
		Order("key_name ASC").
		Find(&metas).Error
	return metas, err
}

func (r *SecretMetadataRepository) FindByKey(projectID, envID uuid.UUID, keyName string) (*models.SecretMetadata, error) {
	var meta models.SecretMetadata
	err := r.db.First(&meta, "project_id = ? AND environment_id = ? AND key_name = ?",
		projectID, envID, keyName).Error
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

func (r *SecretMetadataRepository) DeleteByKey(projectID, envID uuid.UUID, keyName string) error {
	return r.db.Where("project_id = ? AND environment_id = ? AND key_name = ?",
		projectID, envID, keyName).
		Delete(&models.SecretMetadata{}).Error
}

// IncrementVersion bumps vault_version and last_modified_at for a key.
func (r *SecretMetadataRepository) IncrementVersion(projectID, envID uuid.UUID, keyName string, userID uuid.UUID) error {
	return r.db.Model(&models.SecretMetadata{}).
		Where("project_id = ? AND environment_id = ? AND key_name = ?", projectID, envID, keyName).
		Updates(map[string]interface{}{
			"vault_version":    gorm.Expr("vault_version + 1"),
			"last_modified_at": time.Now(),
			"created_by_id":    userID,
		}).Error
}
