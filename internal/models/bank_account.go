package models

import (
	"time"

	"github.com/google/uuid"
)

type BankAccount struct {
	BankAccountID    uuid.UUID  `json:"bank_account_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	UserID           *uuid.UUID `json:"user_id"`
	Name             *string    `json:"name"`
	PlaidAccessToken *string    `json:"plaid_access_token"`
	PlaidItemID      *string    `json:"plaid_item_id"`
	CreatedAt        time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
