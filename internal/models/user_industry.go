package models

import (
	"time"

	"github.com/google/uuid"
)

type UserIndustry struct {
	UserIndustryID uuid.UUID `json:"user_industry_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	IndustryID     uuid.UUID `json:"industry_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	UserID         uuid.UUID `json:"user_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt      time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
