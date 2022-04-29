package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID         uuid.UUID  `json:"order_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	ParentOrderID   *uuid.UUID `json:"parent_order_id" gorm:"type:uuid"`
	MatchingOrderID *uuid.UUID `json:"matching_order_id" gorm:"type:uuid"`
	UserID          *uuid.UUID `json:"user_id" gorm:"type:uuid"`
	PortfolioID     *uuid.UUID `json:"portfolio_id" gorm:"type:uuid"`
	Amount          int        `json:"amount"`
	Side            string     `json:"side"`
	Status          string     `json:"status"`
	Symbol          *string    `json:"symbol"`
	AlpacaOrderID   *string    `json:"alpaca_order_id"`
	CompletedAt     time.Time  `json:"completed_at"`
	CreatedAt       time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	ChildOrders     []Order    `json:"child_orders" gorm:"foreignKey:OrderID"`
}
