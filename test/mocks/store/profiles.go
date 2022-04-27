package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetProfile(profileID uuid.UUID) (*models.Profile, error) {
	args := mockStore.Called(profileID)

	return args.Get(0).(*models.Profile), args.Error(1)
}

func (mockStore *MockStore) ListProfiles(query map[string]interface{}) ([]models.Profile, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.Profile), args.Error(1)
}

func (mockStore *MockStore) CreateProfile(profilePayload models.Profile) (*models.Profile, error) {
	args := mockStore.Called(profilePayload)

	return args.Get(0).(*models.Profile), args.Error(1)
}

func (mockStore *MockStore) ModifyProfile(
	profileID uuid.UUID,
	profilePayload models.Profile,
) (*models.Profile, error) {
	args := mockStore.Called(profileID, profilePayload)

	return args.Get(0).(*models.Profile), args.Error(1)
}

func (mockStore *MockStore) DeleteProfile(profileID uuid.UUID) error {
	args := mockStore.Called(profileID)

	return args.Error(0)
}
