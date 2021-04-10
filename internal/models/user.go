package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	FirstName string    `json:"first_name"`
}
