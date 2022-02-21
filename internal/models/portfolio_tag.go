package models

import (
	"time"

	"github.com/google/uuid"
)

type PortfolioTag struct {
	PortfolioTagID uuid.UUID `json:"portfolio_tag_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	TagID          uuid.UUID `json:"tag_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	PortfolioID    uuid.UUID `json:"portfolio_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt      time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
