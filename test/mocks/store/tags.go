package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetTag(tagID uuid.UUID) (*models.Tag, error) {
	args := mockStore.Called(tagID)

	return args.Get(0).(*models.Tag), args.Error(1)
}

func (mockStore *MockStore) ListTags(query map[string]interface{}) ([]models.Tag, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.Tag), args.Error(1)
}

func (mockStore *MockStore) CreateTag(tagPayload models.Tag) (*models.Tag, error) {
	args := mockStore.Called(tagPayload)

	return args.Get(0).(*models.Tag), args.Error(1)
}

func (mockStore *MockStore) ModifyTag(
	tagID uuid.UUID,
	tagPayload models.Tag,
) (*models.Tag, error) {
	args := mockStore.Called(tagID, tagPayload)

	return args.Get(0).(*models.Tag), args.Error(1)
}

func (mockStore *MockStore) DeleteTag(tagID uuid.UUID) error {
	args := mockStore.Called(tagID)

	return args.Error(0)
}
