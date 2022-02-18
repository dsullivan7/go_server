package models

import (
	"time"

	"github.com/google/uuid"
)

type Security struct {
	SecurityID           uuid.UUID  `json:"security_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Symbol                  string    `json:"symbol"`
	Name                  string    `json:"name"`
	Beta                  string    `json:"beta"`
	CreatedAt               time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt               time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
