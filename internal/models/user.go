package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `json:"user_id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Auth0ID   *string   `json:"auth0_id"`
	DwollaCustomerID   *string   `json:"dwolla_customer_id"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
	PhoneNumber  *string   `json:"phone_number"`
	Email  *string   `json:"email"`
	Address1  *string   `json:"address_1"`
	City  *string   `json:"city"`
	State  *string   `json:"state"`
	PostalCode  *string   `json:"postal_code"`
	DateOfBirth  *string   `json:"-"`
	SSN  *string   `json:"-"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
