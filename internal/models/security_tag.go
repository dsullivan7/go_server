package models

import (
	"time"

	"github.com/google/uuid"
)

type SecurityTag struct {
	SecurityTagID uuid.UUID `json:"security_tag_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	TagID         uuid.UUID `json:"tag_id" gorm:"type:uuid"`
	SecurityID    uuid.UUID `json:"security_id" gorm:"type:uuid"`
	CreatedAt     time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
