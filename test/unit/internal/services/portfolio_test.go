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

	securityID1 := uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000001"))
	securityID2 := uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000002"))
	securityID3 := uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000003"))
	securityID4 := uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000004"))
	securityID5 := uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000005"))
	securityID6 := uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000006"))

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
				services.PortfolioHolding{Symbol: "symbol1", Name: "", Amount: .5},
				services.PortfolioHolding{Symbol: "symbol2", Name: "", Amount: .5},
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
				services.PortfolioHolding{Symbol: "symbol1", Name: "", Amount: .6667},
				services.PortfolioHolding{Symbol: "symbol2", Name: "", Amount: 0.3333},
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
				services.PortfolioHolding{Symbol: "symbol1", Name: "", Amount: .25},
				services.PortfolioHolding{Symbol: "symbol2", Name: "", Amount: .50},
				services.PortfolioHolding{Symbol: "symbol3", Name: "", Amount: .25},
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
				services.PortfolioHolding{Symbol: "symbol1", Name: "", Amount: .50},
				services.PortfolioHolding{Symbol: "symbol2", Name: "", Amount: .50},
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
				services.PortfolioHolding{Symbol: "symbol1", Name: "", Amount: .50},
				services.PortfolioHolding{Symbol: "symbol2", Name: "", Amount: .50},
			},
		},
		{
			name: "rounding",
			portfolio: models.Portfolio{
				Risk: 3,
			},
			portfolioTags: []models.PortfolioTag{
				models.PortfolioTag{TagID: tagID1},
			},
			securities: []models.Security{
				models.Security{Symbol: "symbol1", SecurityID: securityID1},
				models.Security{Symbol: "symbol2", SecurityID: securityID2},
				models.Security{Symbol: "symbol3", SecurityID: securityID3},
				models.Security{Symbol: "symbol4", SecurityID: securityID4},
				models.Security{Symbol: "symbol5", SecurityID: securityID5},
				models.Security{Symbol: "symbol6", SecurityID: securityID6},
			},
			securityTags: []models.SecurityTag{
				models.SecurityTag{SecurityID: securityID1, TagID: tagID1},
				models.SecurityTag{SecurityID: securityID2, TagID: tagID1},
				models.SecurityTag{SecurityID: securityID3, TagID: tagID1},
				models.SecurityTag{SecurityID: securityID4, TagID: tagID1},
				models.SecurityTag{SecurityID: securityID5, TagID: tagID1},
				models.SecurityTag{SecurityID: securityID6, TagID: tagID1},
			},
			target: []services.PortfolioHolding{
				services.PortfolioHolding{Symbol: "symbol1", Name: "", Amount: .1667},
				services.PortfolioHolding{Symbol: "symbol2", Name: "", Amount: .1667},
				services.PortfolioHolding{Symbol: "symbol3", Name: "", Amount: .1667},
				services.PortfolioHolding{Symbol: "symbol4", Name: "", Amount: .1667},
				services.PortfolioHolding{Symbol: "symbol5", Name: "", Amount: .1667},
				services.PortfolioHolding{Symbol: "symbol6", Name: "", Amount: .1665},
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
