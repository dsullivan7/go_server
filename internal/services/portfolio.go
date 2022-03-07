package services

import (
	"go_server/internal/models"
	"math"
)

type IService interface {
	ListPortfolioHoldings(
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
	Symbol string  `json:"symbol"`
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

const portfolioTotal = 1.0
const roundValue = 100

// ListPortfolioHoldings retreives a set of portfolio holdings
// according to the specified portfolio and portfolio tags.
func (srvc *Service) ListPortfolioHoldings(
	portfolio models.Portfolio,
	portfolioTags []models.PortfolioTag,
	securities []models.Security,
	securityTags []models.SecurityTag,
) []PortfolioHolding {
	// create a lookup map for securities
	securityMap := map[string]models.Security{}
	for _, security := range securities {
		securityMap[security.SecurityID.String()] = security
	}

	// create a map of securities to the tag weights
	securityWeightMap := map[string]int{}
	totalWeight := 0

	for _, portfolioTag := range portfolioTags {
		for _, securityTag := range securityTags {
			if securityTag.TagID == portfolioTag.TagID {
				// look up the security
				security := securityMap[securityTag.SecurityID.String()]
				securityWeightMap[security.SecurityID.String()]++
				totalWeight++
			}
		}
	}

	portfolioHoldings := make([]PortfolioHolding, len(securityWeightMap))
	currentIndex := 0
	remaining := portfolioTotal

	for securityID, securityWeight := range securityWeightMap {
		var amount float64
		if currentIndex+1 == len(securityWeightMap) {
			amount = remaining
		} else {
			raw := (float64(securityWeight) / float64(totalWeight)) * float64(portfolioTotal)
			// round the amount to 2 decimal places
			amount = math.Round(raw*roundValue) / roundValue
			remaining -= amount
		}

		portfolioHoldings[currentIndex] = PortfolioHolding{
			Symbol: securityMap[securityID].Symbol,
			Name:   securityMap[securityID].Name,
			Amount: amount,
		}
		currentIndex++
	}

	return portfolioHoldings
}
