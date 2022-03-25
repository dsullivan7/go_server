package models

import (
	"time"

	"github.com/google/uuid"
)

type Credential struct {
	CredentialID uuid.UUID `json:"credential_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
