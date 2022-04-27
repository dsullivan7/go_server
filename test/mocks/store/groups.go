package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetGroup(groupID uuid.UUID) (*models.Group, error) {
	args := mockStore.Called(groupID)

	return args.Get(0).(*models.Group), args.Error(1)
}

func (mockStore *MockStore) ListGroups(query map[string]interface{}) ([]models.Group, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.Group), args.Error(1)
}

func (mockStore *MockStore) CreateGroup(groupPayload models.Group) (*models.Group, error) {
	args := mockStore.Called(groupPayload)

	return args.Get(0).(*models.Group), args.Error(1)
}

func (mockStore *MockStore) ModifyGroup(
	groupID uuid.UUID,
	groupPayload models.Group,
) (*models.Group, error) {
	args := mockStore.Called(groupID, groupPayload)

	return args.Get(0).(*models.Group), args.Error(1)
}

func (mockStore *MockStore) DeleteGroup(groupID uuid.UUID) error {
	args := mockStore.Called(groupID)

	return args.Error(0)
}
