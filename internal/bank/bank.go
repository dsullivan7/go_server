package bank

import (
	"go_server/internal/models"
)

type Bank interface {
	CreateTransfer(source models.BankAccount, destination models.BankAccount) models.BankTransfer
}
