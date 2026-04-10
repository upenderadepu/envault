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
func (r *UserRepository) FindOrCreate(supabaseUID, email string) (*models.User, error) {
	var user models.User
	err := r.db.Where(models.User{SupabaseUID: supabaseUID}).
		Attrs(models.User{Email: email}).
		FirstOrCreate(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
