package models

import (
	"time"

	"github.com/google/uuid"
)

type Industry struct {
	IndustryID uuid.UUID `json:"industry_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name       *string   `json:"name"`
	CreatedAt  time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
