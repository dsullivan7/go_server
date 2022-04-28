package models

import (
	"time"

	"github.com/google/uuid"
)

type BankAccount struct {
	BankAccountID           uuid.UUID  `json:"bank_account_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	UserID                  *uuid.UUID `json:"user_id"`
	Name                    *string    `json:"name"`
	PlaidAccessToken        *string    `json:"plaid_access_token"`
	PlaidAccountID          *string    `json:"plaid_account_id"`
	DwollaFundingSourceID          *string    `json:"dwolla_funding_source_id"`
	AlpacaACHRelationshipID *string    `json:"alpaca_ach_relationship_id"`
	CreatedAt               time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt               time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
