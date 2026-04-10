package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	SupabaseUID string    `gorm:"column:supabase_uid;uniqueIndex;not null" json:"supabase_uid"`
	Email       string    `gorm:"uniqueIndex;not null" json:"email"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
}
