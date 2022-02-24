package services_test

import (
	"go_server/internal/models"
	"go_server/internal/services"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPortfolio(tParent *testing.T) {
	tParent.Parallel()

	srvc := services.NewService()

	type testCase struct {
		name          string
		portfolio     models.Portfolio
		portfolioTags []models.PortfolioTag
		securities    []models.Security
		securityTags  []models.SecurityTag
		target        []services.PortfolioHolding
	}

	tagID1 := uuid.New()
	tagID2 := uuid.New()
	tagID3 := uuid.New()

	securityID1 := uuid.New()
	securityID2 := uuid.New()
	securityID3 := uuid.New()

	tests := []testCase{
		{
			name: "simple",
			portfolio: models.Portfolio{
				Risk: 3,
			},
			portfolioTags: []models.PortfolioTag{
				models.PortfolioTag{TagID: tagID1},
				models.PortfolioTag{TagID: tagID2},
			},
			securities: []models.Security{
				models.Security{Symbol: "symbol1", SecurityID: securityID1},
				models.Security{Symbol: "symbol2", SecurityID: securityID2},
			},
			securityTags: []models.SecurityTag{
				models.SecurityTag{SecurityID: securityID1, TagID: tagID1},
				models.SecurityTag{SecurityID: securityID2, TagID: tagID2},
			},
			target: []services.PortfolioHolding{
				services.PortfolioHolding{Symbol: "symbol1", Amount: 50},
				services.PortfolioHolding{Symbol: "symbol2", Amount: 50},
			},
		},
		{
			name: "simple_multiple_tags_1",
			portfolio: models.Portfolio{
				Risk: 3,
			},
			portfolioTags: []models.PortfolioTag{
				models.PortfolioTag{TagID: tagID1},
				models.PortfolioTag{TagID: tagID2},
			},
			securities: []models.Security{
				models.Security{Symbol: "symbol1", SecurityID: securityID1},
				models.Security{Symbol: "symbol2", SecurityID: securityID2},
			},
			securityTags: []models.SecurityTag{
				models.SecurityTag{SecurityID: securityID1, TagID: tagID1},
				models.SecurityTag{SecurityID: securityID1, TagID: tagID2},
				models.SecurityTag{SecurityID: securityID2, TagID: tagID2},
			},
			target: []services.PortfolioHolding{
				services.PortfolioHolding{Symbol: "symbol1", Amount: 66.67},
				services.PortfolioHolding{Symbol: "symbol2", Amount: 33.33},
			},
		},
		{
			name: "simple_multiple_tags_2",
			portfolio: models.Portfolio{
				Risk: 3,
			},
			portfolioTags: []models.PortfolioTag{
				models.PortfolioTag{TagID: tagID1},
				models.PortfolioTag{TagID: tagID2},
				models.PortfolioTag{TagID: tagID3},
			},
			securities: []models.Security{
				models.Security{Symbol: "symbol1", SecurityID: securityID1},
				models.Security{Symbol: "symbol2", SecurityID: securityID2},
				models.Security{Symbol: "symbol3", SecurityID: securityID3},
			},
			securityTags: []models.SecurityTag{
				models.SecurityTag{SecurityID: securityID1, TagID: tagID1},
				models.SecurityTag{SecurityID: securityID2, TagID: tagID1},
				models.SecurityTag{SecurityID: securityID2, TagID: tagID2},
				models.SecurityTag{SecurityID: securityID3, TagID: tagID3},
			},
			target: []services.PortfolioHolding{
				services.PortfolioHolding{Symbol: "symbol1", Amount: 25},
				services.PortfolioHolding{Symbol: "symbol2", Amount: 50},
				services.PortfolioHolding{Symbol: "symbol3", Amount: 25},
			},
		},
		{
			name: "simple_no_tag_match",
			portfolio: models.Portfolio{
				Risk: 3,
			},
			portfolioTags: []models.PortfolioTag{
				models.PortfolioTag{TagID: tagID1},
				models.PortfolioTag{TagID: tagID2},
			},
			securities: []models.Security{
				models.Security{Symbol: "symbol1", SecurityID: securityID1},
				models.Security{Symbol: "symbol2", SecurityID: securityID2},
				models.Security{Symbol: "symbol3", SecurityID: securityID3},
			},
			securityTags: []models.SecurityTag{
				models.SecurityTag{SecurityID: securityID1, TagID: tagID1},
				models.SecurityTag{SecurityID: securityID2, TagID: tagID2},
				models.SecurityTag{SecurityID: securityID3, TagID: tagID3},
			},
			target: []services.PortfolioHolding{
				services.PortfolioHolding{Symbol: "symbol1", Amount: 50},
				services.PortfolioHolding{Symbol: "symbol2", Amount: 50},
			},
		},
		{
			name: "multi_symbol_multi_matching",
			portfolio: models.Portfolio{
				Risk: 3,
			},
			portfolioTags: []models.PortfolioTag{
				models.PortfolioTag{TagID: tagID1},
				models.PortfolioTag{TagID: tagID2},
			},
			securities: []models.Security{
				models.Security{Symbol: "symbol1", SecurityID: securityID1},
				models.Security{Symbol: "symbol2", SecurityID: securityID2},
			},
			securityTags: []models.SecurityTag{
				models.SecurityTag{SecurityID: securityID1, TagID: tagID1},
				models.SecurityTag{SecurityID: securityID1, TagID: tagID2},
				models.SecurityTag{SecurityID: securityID2, TagID: tagID1},
				models.SecurityTag{SecurityID: securityID2, TagID: tagID2},
			},
			target: []services.PortfolioHolding{
				services.PortfolioHolding{Symbol: "symbol1", Amount: 50},
				services.PortfolioHolding{Symbol: "symbol2", Amount: 50},
			},
		},
	}

	for _, testCase := range tests {
		tc := testCase
		tParent.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := srvc.ListPortfolioHoldings(tc.portfolio, tc.portfolioTags, tc.securities, tc.securityTags)
			assert.ElementsMatch(t, tc.target, actual)
		})
	}
}
