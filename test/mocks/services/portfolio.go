package services

import (
	"go_server/internal/models"
	"go_server/internal/services"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func NewMockService() *MockService {
	return &MockService{}
}

func (mockService *MockService) GetPortfolioHoldings(
	portfolio models.Portfolio,
	portfolioTags []models.PortfolioTag,
	securities []models.Security,
	securityTags []models.SecurityTag,
) []services.PortfolioHolding {
	args := mockService.Called(portfolio, portfolioTags, securities, securityTags)

	return args.Get(0).([]services.PortfolioHolding)
}
