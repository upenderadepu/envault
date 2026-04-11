package models

import (
	"time"

	"github.com/google/uuid"
)

type TeamMember struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProjectID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"project_id"`
	UserID             uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	User               *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role               string     `gorm:"not null" json:"role"` // admin, developer, ci
	VaultPolicyName    string     `json:"vault_policy_name,omitempty"`
	VaultTokenAccessor string     `json:"-"`                                          // never expose in API responses
	InviteCode         string     `gorm:"uniqueIndex" json:"-"`                       // short code for invite-by-share
	IsActive           bool       `gorm:"not null;default:true" json:"is_active"`
	InvitedAt          time.Time  `gorm:"not null;default:now()" json:"invited_at"`
	JoinedAt           *time.Time `json:"joined_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
