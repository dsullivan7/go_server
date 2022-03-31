package models

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ItemID    uuid.UUID `json:"item_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Price     float64   `json:"price"`
	Name      string    `json:"name"`
	InvoiceID uuid.UUID `json:"invoice_id" gorm:"type:uuid"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
