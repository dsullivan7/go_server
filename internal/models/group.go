package models

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	GroupID         uuid.UUID `json:"group_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name            string    `json:"name"`
	APIClientKey    string    `json:"api_client_key"`
	APIClientSecret string    `json:"api_client_secret"`
	CreatedAt       time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
