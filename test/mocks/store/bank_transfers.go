package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetBankTransfer(bankTransferID uuid.UUID) (*models.BankTransfer, error) {
	args := mockStore.Called(bankTransferID)

	return args.Get(0).(*models.BankTransfer), args.Error(1)
}

func (mockStore *MockStore) ListBankTransfers(query map[string]interface{}) ([]models.BankTransfer, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.BankTransfer), args.Error(1)
}

func (mockStore *MockStore) CreateBankTransfer(bankTransferPayload models.BankTransfer) (*models.BankTransfer, error) {
	args := mockStore.Called(bankTransferPayload)

	return args.Get(0).(*models.BankTransfer), args.Error(1)
}

func (mockStore *MockStore) ModifyBankTransfer(
	bankTransferID uuid.UUID,
	bankTransferPayload models.BankTransfer,
) (*models.BankTransfer, error) {
	args := mockStore.Called(bankTransferID, bankTransferPayload)

	return args.Get(0).(*models.BankTransfer), args.Error(1)
}

func (mockStore *MockStore) DeleteBankTransfer(bankTransferID uuid.UUID) error {
	args := mockStore.Called(bankTransferID)

	return args.Error(0)
}
