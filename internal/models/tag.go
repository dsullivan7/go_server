package models

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	TagID     uuid.UUID `json:"tag_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name      *string   `json:"name"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
