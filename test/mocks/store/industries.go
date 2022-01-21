package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetIndustry(industryID uuid.UUID) (*models.Industry, error) {
	args := mockStore.Called(industryID)

	return args.Get(0).(*models.Industry), args.Error(1)
}

func (mockStore *MockStore) ListIndustries(query map[string]interface{}) ([]models.Industry, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.Industry), args.Error(1)
}

func (mockStore *MockStore) CreateIndustry(industryPayload models.Industry) (*models.Industry, error) {
	args := mockStore.Called(industryPayload)

	return args.Get(0).(*models.Industry), args.Error(1)
}

func (mockStore *MockStore) ModifyIndustry(
	industryID uuid.UUID,
	industryPayload models.Industry,
) (*models.Industry, error) {
	args := mockStore.Called(industryID, industryPayload)

	return args.Get(0).(*models.Industry), args.Error(1)
}

func (mockStore *MockStore) DeleteIndustry(industryID uuid.UUID) error {
	args := mockStore.Called(industryID)

	return args.Error(0)
}
