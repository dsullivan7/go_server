package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetReview(reviewID uuid.UUID) (*models.Review, error) {
	args := mockStore.Called(reviewID)

	return args.Get(0).(*models.Review), args.Error(1)
}

func (mockStore *MockStore) ListReviews(query map[string]interface{}) ([]models.Review, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.Review), args.Error(1)
}

func (mockStore *MockStore) CreateReview(reviewPayload models.Review) (*models.Review, error) {
	args := mockStore.Called(reviewPayload)

	return args.Get(0).(*models.Review), args.Error(1)
}

func (mockStore *MockStore) ModifyReview(reviewID uuid.UUID, reviewPayload models.Review) (*models.Review, error) {
	args := mockStore.Called(reviewID, reviewPayload)

	return args.Get(0).(*models.Review), args.Error(1)
}

func (mockStore *MockStore) DeleteReview(reviewID uuid.UUID) error {
	args := mockStore.Called(reviewID)

	return args.Error(0)
}
