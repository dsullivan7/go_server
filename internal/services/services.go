package services

import (
	"time"

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
	GetBalance(orders []models.Order, interest float64, currentTime time.Time) (int, int, int)
}

type Service struct {
}

func NewService() IService {
	return &Service{}
}
