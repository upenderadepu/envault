package repository

import (
	"github.com/bhartiyaanshul/envault/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindBySupabaseUID(uid string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "supabase_uid = ?", uid).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindOrCreate looks up a user by SupabaseUID, creating one if not found.
// If the user exists but their email changed, it updates the email.
func (r *UserRepository) FindOrCreate(supabaseUID, email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("supabase_uid = ?", supabaseUID).First(&user).Error
	if err == nil {
		// User exists — update email if changed
		if user.Email != email && email != "" {
			r.db.Model(&user).Update("email", email)
		}
		return &user, nil
	}

	// User not found by UID — check if email already exists (different auth provider)
	var existingByEmail models.User
	if err := r.db.Where("email = ?", email).First(&existingByEmail).Error; err == nil {
		// Update the existing record with the new supabase_uid
		r.db.Model(&existingByEmail).Update("supabase_uid", supabaseUID)
		return &existingByEmail, nil
	}

	// Completely new user — create
	user = models.User{
		SupabaseUID: supabaseUID,
		Email:       email,
	}
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
