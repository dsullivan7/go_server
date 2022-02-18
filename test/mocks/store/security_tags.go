package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetSecurityTag(securityTagID uuid.UUID) (*models.SecurityTag, error) {
	args := mockStore.Called(securityTagID)

	return args.Get(0).(*models.SecurityTag), args.Error(1)
}

func (mockStore *MockStore) ListSecurityTags(query map[string]interface{}) ([]models.SecurityTag, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.SecurityTag), args.Error(1)
}

func (mockStore *MockStore) CreateSecurityTag(
	securityTagPayload models.SecurityTag,
) (*models.SecurityTag, error) {
	args := mockStore.Called(securityTagPayload)

	return args.Get(0).(*models.SecurityTag), args.Error(1)
}

func (mockStore *MockStore) ModifySecurityTag(
	securityTagID uuid.UUID,
	securityTagPayload models.SecurityTag,
) (*models.SecurityTag, error) {
	args := mockStore.Called(securityTagID, securityTagPayload)

	return args.Get(0).(*models.SecurityTag), args.Error(1)
}

func (mockStore *MockStore) DeleteSecurityTag(securityTagID uuid.UUID) error {
	args := mockStore.Called(securityTagID)

	return args.Error(0)
}
