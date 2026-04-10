package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Action constants for audit logging.
const (
	ActionProjectCreate     = "project.create"
	ActionProjectDelete     = "project.delete"
	ActionSecretRead        = "secret.read"
	ActionSecretWrite       = "secret.write"
	ActionSecretDelete      = "secret.delete"
	ActionMemberInvite      = "member.invite"
	ActionMemberRemove      = "member.remove"
	ActionCredentialsRotate = "credentials.rotate"
)

// AuditLog is append-only. The database enforces immutability via
// REVOKE UPDATE, DELETE on the audit_logs table.
type AuditLog struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProjectID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	UserID       *uuid.UUID     `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Action       string         `gorm:"not null" json:"action"`
	ResourcePath string         `gorm:"not null" json:"resource_path"`
	IPAddress    string         `json:"ip_address,omitempty"`
	UserAgent    string         `json:"user_agent,omitempty"`
	RequestID    string         `json:"request_id,omitempty"`
	Metadata     datatypes.JSON `gorm:"default:'{}'" json:"metadata"`
	CreatedAt    time.Time      `gorm:"not null;default:now()" json:"created_at"`
}
