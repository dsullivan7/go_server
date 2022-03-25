package store

import (
	"go_server/internal/models"
)

func (mockStore *MockStore) CreateCredential(tagPayload models.Credential) (*models.Credential, error) {
	args := mockStore.Called(tagPayload)

	return args.Get(0).(*models.Credential), args.Error(1)
}
