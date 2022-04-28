package bank

import (
	"go_server/internal/models"
)

type Bank interface {
	CreateTransfer(source models.BankAccount, destination models.BankAccount, amount int) (*models.BankTransfer, error)
	CreateCustomer(user models.User) (*models.User, error)
	CreateBank(user models.User, plaidProcessorToken string) (*models.BankAccount, error)
}
