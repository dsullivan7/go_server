package models

import (
	"time"

	"github.com/google/uuid"
)

type Portfolio struct {
	PortfolioID   uuid.UUID  `json:"portfolio_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	UserID *uuid.UUID `json:"user_id"`
	Risk  int  `json:"risk"`
	CreatedAt  time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
