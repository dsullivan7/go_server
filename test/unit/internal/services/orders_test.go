package services_test

import (
	"go_server/internal/models"
	"go_server/internal/services"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestOrders(tParent *testing.T) {
	tParent.Parallel()

	srvc := services.NewService()

	type testCase struct {
		name          string
		openOrders     []models.Order
		netSecurityValue int
		netCashValue int
		target        []models.Order
	}

	uuid1 := uuid.New()

	tests := []testCase{
		{
			name: "simple",
			openOrders: []models.Order{
				models.Order{
					OrderID: uuid1,
					Side: "buy",
					Amount: 200,
					ChildOrders: []models.Order{
						models.Order{ Amount: 100, Side: "buy" },
					},
				},
			},
			netSecurityValue: 100,
			netCashValue: 0,
			target: []models.Order{
				models.Order{
					ParentOrderID: &uuid1,
					Amount: 100,
					Side: "buy",
				},
			},
		},
	}

	for _, testCase := range tests {
		tc := testCase
		tParent.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := srvc.GetOrders(tc.openOrders, tc.netSecurityValue, tc.netCashValue)
			assert.ElementsMatch(t, tc.target, actual)
		})
	}
}
