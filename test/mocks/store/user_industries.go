package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetUserIndustry(userIndustryID uuid.UUID) (*models.UserIndustry, error) {
	args := mockStore.Called(userIndustryID)

	return args.Get(0).(*models.UserIndustry), args.Error(1)
}

func (mockStore *MockStore) ListUserIndustries(query map[string]interface{}) ([]models.UserIndustry, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.UserIndustry), args.Error(1)
}

func (mockStore *MockStore) CreateUserIndustry(userIndustryPayload models.UserIndustry) (*models.UserIndustry, error) {
	args := mockStore.Called(userIndustryPayload)

	return args.Get(0).(*models.UserIndustry), args.Error(1)
}

func (mockStore *MockStore) ModifyUserIndustry(userIndustryID uuid.UUID, userIndustryPayload models.UserIndustry) (*models.UserIndustry, error) {
	args := mockStore.Called(userIndustryID, userIndustryPayload)

	return args.Get(0).(*models.UserIndustry), args.Error(1)
}

func (mockStore *MockStore) DeleteUserIndustry(userIndustryID uuid.UUID) error {
	args := mockStore.Called(userIndustryID)

	return args.Error(0)
}
