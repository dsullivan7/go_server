package services

import (
	"go_server/internal/models"
)

type IService interface {
	GetPortfolio(
		models.Portfolio,
		[]models.PortfolioTag,
		[]models.Security,
		[]models.SecurityTag,
	) []PortfolioHolding
}

type Service struct {
}

func NewService() IService {
	return &Service{}
}

type PortfolioHolding struct {
	Symbol string
	Amount float64
}

func (srvc *Service) GetPortfolio(
	models.Portfolio,
	[]models.PortfolioTag,
	[]models.Security,
	[]models.SecurityTag,
) []PortfolioHolding {
	return []PortfolioHolding{}
}
