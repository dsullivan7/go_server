package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetPortfolioTag(portfolioTagID uuid.UUID) (*models.PortfolioTag, error) {
	args := mockStore.Called(portfolioTagID)

	return args.Get(0).(*models.PortfolioTag), args.Error(1)
}

func (mockStore *MockStore) ListPortfolioTags(query map[string]interface{}) ([]models.PortfolioTag, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.PortfolioTag), args.Error(1)
}

func (mockStore *MockStore) CreatePortfolioTag(
	portfolioTagPayload models.PortfolioTag,
) (*models.PortfolioTag, error) {
	args := mockStore.Called(portfolioTagPayload)

	return args.Get(0).(*models.PortfolioTag), args.Error(1)
}

func (mockStore *MockStore) ModifyPortfolioTag(
	portfolioTagID uuid.UUID,
	portfolioTagPayload models.PortfolioTag,
) (*models.PortfolioTag, error) {
	args := mockStore.Called(portfolioTagID, portfolioTagPayload)

	return args.Get(0).(*models.PortfolioTag), args.Error(1)
}

func (mockStore *MockStore) DeletePortfolioTag(portfolioTagID uuid.UUID) error {
	args := mockStore.Called(portfolioTagID)

	return args.Error(0)
}
