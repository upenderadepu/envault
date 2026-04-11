package repository

import (
	"github.com/bhartiyaanshul/envault/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamMemberRepository struct {
	db *gorm.DB
}

func NewTeamMemberRepository(db *gorm.DB) *TeamMemberRepository {
	return &TeamMemberRepository{db: db}
}

func (r *TeamMemberRepository) Create(member *models.TeamMember) error {
	return r.db.Create(member).Error
}

func (r *TeamMemberRepository) FindByProjectAndUser(projectID, userID uuid.UUID) (*models.TeamMember, error) {
	var member models.TeamMember
	err := r.db.First(&member, "project_id = ? AND user_id = ?", projectID, userID).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *TeamMemberRepository) FindByID(id uuid.UUID) (*models.TeamMember, error) {
	var member models.TeamMember
	err := r.db.Preload("User").First(&member, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// FindActiveByProjectID returns active team members with user info preloaded.
func (r *TeamMemberRepository) FindActiveByProjectID(projectID uuid.UUID) ([]models.TeamMember, error) {
	var members []models.TeamMember
	err := r.db.Preload("User").
		Where("project_id = ? AND is_active = ?", projectID, true).
		Find(&members).Error
	return members, err
}

func (r *TeamMemberRepository) Update(member *models.TeamMember) error {
	return r.db.Save(member).Error
}

// FindByAccessor looks up a team member by their Vault token accessor.
func (r *TeamMemberRepository) FindByAccessor(accessor string) (*models.TeamMember, error) {
	var member models.TeamMember
	err := r.db.Preload("User").First(&member, "vault_token_accessor = ? AND is_active = ?", accessor, true).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// FindByInviteCode looks up a pending team member by invite code.
func (r *TeamMemberRepository) FindByInviteCode(code string) (*models.TeamMember, error) {
	var member models.TeamMember
	err := r.db.Preload("User").First(&member, "invite_code = ? AND is_active = ?", code, true).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// GetAccessorsByProjectID returns all non-empty vault token accessors for a project.
func (r *TeamMemberRepository) GetAccessorsByProjectID(projectID uuid.UUID) ([]string, error) {
	var accessors []string
	err := r.db.Model(&models.TeamMember{}).
		Where("project_id = ? AND vault_token_accessor != '' AND vault_token_accessor IS NOT NULL", projectID).
		Pluck("vault_token_accessor", &accessors).Error
	return accessors, err
}
