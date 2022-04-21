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
		openOrders []models.Order,
		childOrders []models.Order,
		netSecurityValue int,
		netCashValue int,
	) []models.Order
}

type Service struct {
}

func NewService() IService {
	return &Service{}
}
