package models

import (
  "time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
  gorm.Model
	UserID    uuid.UUID `json:"user_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	FirstName string    `json:"first_name"`
  CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
  UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
