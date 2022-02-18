package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetSecurity(securityID uuid.UUID) (*models.Security, error) {
	args := mockStore.Called(securityID)

	return args.Get(0).(*models.Security), args.Error(1)
}

func (mockStore *MockStore) ListSecurities(query map[string]interface{}) ([]models.Security, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.Security), args.Error(1)
}

func (mockStore *MockStore) CreateSecurity(securityPayload models.Security) (*models.Security, error) {
	args := mockStore.Called(securityPayload)

	return args.Get(0).(*models.Security), args.Error(1)
}

func (mockStore *MockStore) ModifySecurity(
	securityID uuid.UUID,
	securityPayload models.Security,
) (*models.Security, error) {
	args := mockStore.Called(securityID, securityPayload)

	return args.Get(0).(*models.Security), args.Error(1)
}

func (mockStore *MockStore) DeleteSecurity(securityID uuid.UUID) error {
	args := mockStore.Called(securityID)

	return args.Error(0)
}
