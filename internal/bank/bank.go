package bank

import (
	"go_server/internal/models"
)

type Bank interface {
	CreateTransfer(source models.BankAccount, destination models.BankAccount, amount int) (*models.BankTransfer, error)
	CreateCustomer(user models.User) (*models.User, error)
	CreateBankAccount(user models.User, plaidProcessorToken string) (*models.BankAccount, error)
	CreateWebhook() (*models.Webhook, error)
	GetPlaidAccessor() string
}
