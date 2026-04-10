package repository

import (
	"github.com/bhartiyaanshul/envault/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EnvironmentRepository struct {
	db *gorm.DB
}

func NewEnvironmentRepository(db *gorm.DB) *EnvironmentRepository {
	return &EnvironmentRepository{db: db}
}

func (r *EnvironmentRepository) CreateBatch(envs []models.Environment) error {
	return r.db.Create(&envs).Error
}

func (r *EnvironmentRepository) FindByProjectID(projectID uuid.UUID) ([]models.Environment, error) {
	var envs []models.Environment
	err := r.db.Where("project_id = ?", projectID).Find(&envs).Error
	return envs, err
}

func (r *EnvironmentRepository) FindByProjectAndName(projectID uuid.UUID, name string) (*models.Environment, error) {
	var env models.Environment
	err := r.db.First(&env, "project_id = ? AND name = ?", projectID, name).Error
	if err != nil {
		return nil, err
	}
	return &env, nil
}
