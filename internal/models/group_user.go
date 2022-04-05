package models

import (
	"time"

	"github.com/google/uuid"
)

type GroupUser struct {
	GroupUserID    uuid.UUID `json:"group_user_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	GroupID          uuid.UUID `json:"group_id" gorm:"type:uuid"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
