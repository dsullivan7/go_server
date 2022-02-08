package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetBrokerageAccount(brokerageAccountID uuid.UUID) (*models.BrokerageAccount, error) {
	args := mockStore.Called(brokerageAccountID)

	return args.Get(0).(*models.BrokerageAccount), args.Error(1)
}

func (mockStore *MockStore) ListBrokerageAccounts(query map[string]interface{}) ([]models.BrokerageAccount, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.BrokerageAccount), args.Error(1)
}

func (mockStore *MockStore) CreateBrokerageAccount(
	brokerageAccountPayload models.BrokerageAccount,
) (*models.BrokerageAccount, error) {
	args := mockStore.Called(brokerageAccountPayload)

	return args.Get(0).(*models.BrokerageAccount), args.Error(1)
}

func (mockStore *MockStore) ModifyBrokerageAccount(
	brokerageAccountID uuid.UUID,
	brokerageAccountPayload models.BrokerageAccount,
) (*models.BrokerageAccount, error) {
	args := mockStore.Called(brokerageAccountID, brokerageAccountPayload)

	return args.Get(0).(*models.BrokerageAccount), args.Error(1)
}

func (mockStore *MockStore) DeleteBrokerageAccount(brokerageAccountID uuid.UUID) error {
	args := mockStore.Called(brokerageAccountID)

	return args.Error(0)
}
