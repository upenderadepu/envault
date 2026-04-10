package repository

import (
	"github.com/bhartiyaanshul/envault/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

// FindBySlug loads a project with its environments and active team members.
func (r *ProjectRepository) FindBySlug(slug string) (*models.Project, error) {
	var project models.Project
	err := r.db.
		Preload("Environments").
		Preload("TeamMembers", "is_active = ?", true).
		Preload("TeamMembers.User").
		First(&project, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// ListForUser returns all projects where the user is owner or active member.
func (r *ProjectRepository) ListForUser(userID uuid.UUID) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.
		Preload("Environments").
		Where("owner_id = ?", userID).
		Or("id IN (?)",
			r.db.Model(&models.TeamMember{}).
				Select("project_id").
				Where("user_id = ? AND is_active = ?", userID, true),
		).
		Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Project{}, "id = ?", id).Error
}
