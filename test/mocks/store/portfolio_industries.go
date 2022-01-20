package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetPortfolioIndustry(portfolioIndustryID uuid.UUID) (*models.PortfolioIndustry, error) {
	args := mockStore.Called(portfolioIndustryID)

	return args.Get(0).(*models.PortfolioIndustry), args.Error(1)
}

func (mockStore *MockStore) ListPortfolioIndustries(query map[string]interface{}) ([]models.PortfolioIndustry, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.PortfolioIndustry), args.Error(1)
}

func (mockStore *MockStore) CreatePortfolioIndustry(portfolioIndustryPayload models.PortfolioIndustry) (*models.PortfolioIndustry, error) {
	args := mockStore.Called(portfolioIndustryPayload)

	return args.Get(0).(*models.PortfolioIndustry), args.Error(1)
}

func (mockStore *MockStore) ModifyPortfolioIndustry(portfolioIndustryID uuid.UUID, portfolioIndustryPayload models.PortfolioIndustry) (*models.PortfolioIndustry, error) {
	args := mockStore.Called(portfolioIndustryID, portfolioIndustryPayload)

	return args.Get(0).(*models.PortfolioIndustry), args.Error(1)
}

func (mockStore *MockStore) DeletePortfolioIndustry(portfolioIndustryID uuid.UUID) error {
	args := mockStore.Called(portfolioIndustryID)

	return args.Error(0)
}
