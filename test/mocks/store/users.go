package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetUser(userID uuid.UUID) (*models.User, error) {
	args := mockStore.Called(userID)

	return args.Get(0).(*models.User), args.Error(1)
}

func (mockStore *MockStore) ListUsers(query map[string]interface{}) ([]models.User, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.User), args.Error(1)
}

func (mockStore *MockStore) CreateUser(userPayload models.User) (*models.User, error) {
	args := mockStore.Called(userPayload)

	return args.Get(0).(*models.User), args.Error(1)
}

func (mockStore *MockStore) ModifyUser(userID uuid.UUID, userPayload models.User) (*models.User, error) {
	args := mockStore.Called(userID, userPayload)

	return args.Get(0).(*models.User), args.Error(1)
}

func (mockStore *MockStore) DeleteUser(userID uuid.UUID) error {
	args := mockStore.Called(userID)

	return args.Error(0)
}
