package services

import (
	"go_server/internal/models"
)

type IService interface {
	ListPortfolioRecommendations(
		models.Portfolio,
		[]models.PortfolioTag,
		[]models.Security,
		[]models.SecurityTag,
	) []PortfolioHolding
	GetOrders(
		incompleteOrders []models.Order,
		netSecurityValue int,
		netCashValue int,
	) []models.Order
}

type Service struct {
}

func NewService() IService {
	return &Service{}
}
