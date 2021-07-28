package models

import (
  "time"

	"github.com/google/uuid"
)

type Review struct {
	ReviewID    uuid.UUID `json:"review_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	FromUserID  *uuid.UUID `json:"from_user_id"`
	ToUserID  *uuid.UUID `json:"to_user_id"`
	Text *string `json:"text"`
  CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
  UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
