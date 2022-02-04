package models

import (
	"time"

	"github.com/google/uuid"
)

type BrokerageAccount struct {
	BrokerageAccountID uuid.UUID  `json:"brokerage_account_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	UserID             *uuid.UUID `json:"user_id"`
	AlpacaAccountID    *string    `json:"alpaca_account_id"`
	CreatedAt          time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
