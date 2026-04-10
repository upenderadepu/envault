package models

import (
	"time"

	"github.com/google/uuid"
)

type Environment struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProjectID    uuid.UUID `gorm:"type:uuid;not null;index" json:"project_id"`
	Name         string    `gorm:"not null" json:"name"` // development, staging, production
	IsProduction bool      `gorm:"not null;default:false" json:"is_production"`
	CreatedAt    time.Time `json:"created_at"`
}
