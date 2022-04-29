package broker

import (
	"go_server/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockBank struct {
	mock.Mock
}

func NewMockBank() *MockBank {
	return &MockBank{}
}

func (mockBank *MockBank) CreateTransfer(
	source models.BankAccount,
	destination models.BankAccount,
	amount int,
) (*models.BankTransfer, error) {
	args := mockBank.Called(source, destination, amount)

	return args.Get(0).(*models.BankTransfer), args.Error(1)
}

func (mockBank *MockBank) CreateCustomer(user models.User) (*models.User, error) {
	args := mockBank.Called(user)

	return args.Get(0).(*models.User), args.Error(1)
}

func (mockBank *MockBank) CreateBankAccount(user models.User, plaidProcessorToken string) (*models.BankAccount, error) {
	args := mockBank.Called(user, plaidProcessorToken)

	return args.Get(0).(*models.BankAccount), args.Error(1)
}

func (mockBank *MockBank) CreateWebhook() (*models.Webhook, error) {
	args := mockBank.Called()

	return args.Get(0).(*models.Webhook), args.Error(1)
}

func (mockBank *MockBank) GetPlaidAccessor() string {
	args := mockBank.Called()

	return args.String(0)
}
