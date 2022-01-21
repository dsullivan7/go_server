package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetPortfolio(portfolioID uuid.UUID) (*models.Portfolio, error) {
	args := mockStore.Called(portfolioID)

	return args.Get(0).(*models.Portfolio), args.Error(1)
}

func (mockStore *MockStore) ListPortfolios(query map[string]interface{}) ([]models.Portfolio, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.Portfolio), args.Error(1)
}

func (mockStore *MockStore) CreatePortfolio(portfolioPayload models.Portfolio) (*models.Portfolio, error) {
	args := mockStore.Called(portfolioPayload)

	return args.Get(0).(*models.Portfolio), args.Error(1)
}

func (mockStore *MockStore) ModifyPortfolio(
	portfolioID uuid.UUID,
	portfolioPayload models.Portfolio,
) (*models.Portfolio, error) {
	args := mockStore.Called(portfolioID, portfolioPayload)

	return args.Get(0).(*models.Portfolio), args.Error(1)
}

func (mockStore *MockStore) DeletePortfolio(portfolioID uuid.UUID) error {
	args := mockStore.Called(portfolioID)

	return args.Error(0)
}
