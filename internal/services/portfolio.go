package services

import (
	"go_server/internal/models"
	"math"
	"sort"
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
const roundValue = 10000

// ListPortfolioHoldings retreives a set of portfolio holdings
// according to the specified portfolio and portfolio tags.
// nolint:funlen
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

	// sort the securityWeightMap for determinism
	securityIDs := make([]string, 0, len(securityWeightMap))
	for securityID := range securityWeightMap {
		securityIDs = append(securityIDs, securityID)
	}

	sort.Strings(securityIDs)

	for _, securityID := range securityIDs {
		var amount float64
		if currentIndex+1 == len(securityIDs) {
			amount = remaining
		} else {
			amount = (float64(securityWeightMap[securityID]) / float64(totalWeight)) * float64(portfolioTotal)
		}

		// round the amount
		amount = math.Round(amount*roundValue) / roundValue
		remaining -= amount

		println(securityID) //nolint
		println(amount)     //nolint

		portfolioHoldings[currentIndex] = PortfolioHolding{
			Symbol: securityMap[securityID].Symbol,
			Name:   securityMap[securityID].Name,
			Amount: amount,
		}
		currentIndex++
	}

	return portfolioHoldings
}
