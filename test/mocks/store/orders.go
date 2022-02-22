package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetOrder(orderID uuid.UUID) (*models.Order, error) {
	args := mockStore.Called(orderID)

	return args.Get(0).(*models.Order), args.Error(1)
}

func (mockStore *MockStore) ListOrders(query map[string]interface{}) ([]models.Order, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.Order), args.Error(1)
}

func (mockStore *MockStore) CreateOrder(orderPayload models.Order) (*models.Order, error) {
	args := mockStore.Called(orderPayload)

	return args.Get(0).(*models.Order), args.Error(1)
}

func (mockStore *MockStore) ModifyOrder(
	orderID uuid.UUID,
	orderPayload models.Order,
) (*models.Order, error) {
	args := mockStore.Called(orderID, orderPayload)

	return args.Get(0).(*models.Order), args.Error(1)
}

func (mockStore *MockStore) DeleteOrder(orderID uuid.UUID) error {
	args := mockStore.Called(orderID)

	return args.Error(0)
}
