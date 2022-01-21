package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetBankAccount(bankAccountID uuid.UUID) (*models.BankAccount, error) {
	args := mockStore.Called(bankAccountID)

	return args.Get(0).(*models.BankAccount), args.Error(1)
}

func (mockStore *MockStore) ListBankAccounts(query map[string]interface{}) ([]models.BankAccount, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.BankAccount), args.Error(1)
}

func (mockStore *MockStore) CreateBankAccount(bankAccountPayload models.BankAccount) (*models.BankAccount, error) {
	args := mockStore.Called(bankAccountPayload)

	return args.Get(0).(*models.BankAccount), args.Error(1)
}

func (mockStore *MockStore) ModifyBankAccount(
	bankAccountID uuid.UUID,
	bankAccountPayload models.BankAccount,
) (*models.BankAccount, error) {
	args := mockStore.Called(bankAccountID, bankAccountPayload)

	return args.Get(0).(*models.BankAccount), args.Error(1)
}

func (mockStore *MockStore) DeleteBankAccount(bankAccountID uuid.UUID) error {
	args := mockStore.Called(bankAccountID)

	return args.Error(0)
}
