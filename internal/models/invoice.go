package models

import (
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	InvoiceID uuid.UUID  `json:"invoice_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	BankAccountID    *uuid.UUID `json:"bank_account_id" gorm:"type:uuid"`
	CreatedAt   time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
