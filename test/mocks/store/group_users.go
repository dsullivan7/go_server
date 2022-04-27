package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetGroupUser(groupUserID uuid.UUID) (*models.GroupUser, error) {
	args := mockStore.Called(groupUserID)

	return args.Get(0).(*models.GroupUser), args.Error(1)
}

func (mockStore *MockStore) ListGroupUsers(query map[string]interface{}) ([]models.GroupUser, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.GroupUser), args.Error(1)
}

func (mockStore *MockStore) CreateGroupUser(groupUserPayload models.GroupUser) (*models.GroupUser, error) {
	args := mockStore.Called(groupUserPayload)

	return args.Get(0).(*models.GroupUser), args.Error(1)
}

func (mockStore *MockStore) ModifyGroupUser(
	groupUserID uuid.UUID,
	groupUserPayload models.GroupUser,
) (*models.GroupUser, error) {
	args := mockStore.Called(groupUserID, groupUserPayload)

	return args.Get(0).(*models.GroupUser), args.Error(1)
}

func (mockStore *MockStore) DeleteGroupUser(groupUserID uuid.UUID) error {
	args := mockStore.Called(groupUserID)

	return args.Error(0)
}
