package models

import (
	"time"

	"github.com/google/uuid"
)

type BankTransfer struct {
	BankTransferID    uuid.UUID  `json:"bank_transfer_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	UserID           *uuid.UUID `json:"user_id"`
	AlpacaTransferID   *string    `json:"alpaca_transfer_id"`
	CreatedAt        time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
