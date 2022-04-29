package models

import (
	"time"

	"github.com/google/uuid"
)

type Webhook struct {
	WebhookID     uuid.UUID `json:"webhook_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	URL      string   `json:"url"`
	DwollaWebhookID      *string   `json:"dwolla_webhook_id"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
