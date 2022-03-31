package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetItem(itemID uuid.UUID) (*models.Item, error) {
	args := mockStore.Called(itemID)

	return args.Get(0).(*models.Item), args.Error(1)
}

func (mockStore *MockStore) ListItems(query map[string]interface{}) ([]models.Item, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.Item), args.Error(1)
}

func (mockStore *MockStore) CreateItem(itemPayload models.Item) (*models.Item, error) {
	args := mockStore.Called(itemPayload)

	return args.Get(0).(*models.Item), args.Error(1)
}

func (mockStore *MockStore) ModifyItem(
	itemID uuid.UUID,
	itemPayload models.Item,
) (*models.Item, error) {
	args := mockStore.Called(itemID, itemPayload)

	return args.Get(0).(*models.Item), args.Error(1)
}

func (mockStore *MockStore) DeleteItem(itemID uuid.UUID) error {
	args := mockStore.Called(itemID)

	return args.Error(0)
}
