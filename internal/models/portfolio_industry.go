package models

import (
	"time"

	"github.com/google/uuid"
)

type PortfolioIndustry struct {
	PortfolioIndustryID uuid.UUID `json:"portfolio_industry_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	IndustryID     uuid.UUID `json:"industry_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	PortfolioID         uuid.UUID `json:"portfolio_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt      time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
