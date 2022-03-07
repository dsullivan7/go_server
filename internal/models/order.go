package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID       uuid.UUID  `json:"order_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	ParentOrderID *uuid.UUID `json:"parent_order_id" gorm:"type:uuid"`
	UserID        *uuid.UUID `json:"user_id" gorm:"type:uuid"`
	PortfolioID   *uuid.UUID `json:"portfolio_id" gorm:"type:uuid"`
	Amount        float64    `json:"amount"`
	Side          string     `json:"side"`
	Symbol        *string    `json:"symbol"`
	AlpacaOrderID *string    `json:"alpaca_order_id"`
	CreatedAt     time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
