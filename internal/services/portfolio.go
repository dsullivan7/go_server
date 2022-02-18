package services

import (
	"math"
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

const portfolioTotal = 100.0

// GetPortfolio creates a set of portfolio holdings
// according to the specified portfolio and portfolio tags.
func (srvc *Service) GetPortfolio(
	portfolio models.Portfolio,
	portfolioTags []models.PortfolioTag,
	securities []models.Security,
	securityTags []models.SecurityTag,
) []PortfolioHolding {
	// create a lookup map for securities
	securitySymbolMap := map[string]string{}
	for _, security := range securities {
		securitySymbolMap[security.SecurityID.String()] = security.Symbol
	}

	// create a map of securities to the tag weights
	securityWeightMap := map[string]int{}
	totalWeight := 0

	for _, portfolioTag := range portfolioTags {
		for _, securityTag := range securityTags {
			if securityTag.TagID == portfolioTag.TagID {
				// look up the security
				securitySymbol := securitySymbolMap[securityTag.SecurityID.String()]
				securityWeightMap[securitySymbol]++
				totalWeight++
			}
		}
	}

	portfolioHoldings := make([]PortfolioHolding, len(securityWeightMap))
	currentIndex := 0
	remaining := portfolioTotal

	for securitySymbol, securityWeight := range securityWeightMap {
		var amount float64
		if (currentIndex + 1 == len(securityWeightMap)) {
			amount = remaining
		} else {
			raw := (float64(securityWeight) / float64(totalWeight)) * float64(portfolioTotal)
			// round the amount to 2 decimal places
			amount = math.Round(raw * 100) / 100
			remaining -= amount
		}

		portfolioHoldings[currentIndex] = PortfolioHolding{
			Symbol: securitySymbol,
			Amount: amount,
		}
		currentIndex++
	}

	return portfolioHoldings
}