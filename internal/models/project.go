package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name           string         `gorm:"not null" json:"name"`
	Slug           string         `gorm:"uniqueIndex;not null" json:"slug"`
	VaultMountPath string         `gorm:"not null" json:"vault_mount_path"`
	OwnerID        uuid.UUID      `gorm:"type:uuid;not null" json:"owner_id"`
	Owner          *User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Environments   []Environment  `gorm:"foreignKey:ProjectID" json:"environments,omitempty"`
	TeamMembers    []TeamMember   `gorm:"foreignKey:ProjectID" json:"team_members,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
