package models

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ProfileID uuid.UUID `json:"profile_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Username     string
	Password     string
	Type     string `json:"type"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
