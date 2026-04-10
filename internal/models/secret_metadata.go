package models

import (
	"time"

	"github.com/google/uuid"
)

// SecretMetadata tracks which keys exist and their versions.
// Secret VALUES are never stored here — they live only in Vault.
type SecretMetadata struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProjectID      uuid.UUID    `gorm:"type:uuid;not null" json:"project_id"`
	EnvironmentID  uuid.UUID    `gorm:"type:uuid;not null" json:"environment_id"`
	Environment    *Environment `gorm:"foreignKey:EnvironmentID" json:"environment,omitempty"`
	KeyName        string       `gorm:"not null" json:"key_name"`
	VaultPath      string       `gorm:"not null" json:"vault_path"`
	CreatedByID    uuid.UUID    `gorm:"type:uuid;not null" json:"created_by_id"`
	VaultVersion   int          `gorm:"not null;default:1" json:"vault_version"`
	LastModifiedAt time.Time    `gorm:"not null;default:now()" json:"last_modified_at"`
	CreatedAt      time.Time    `json:"created_at"`
}
