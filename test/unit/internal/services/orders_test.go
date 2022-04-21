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

	uuid1 := uuid.New()
	uuid2 := uuid.New()

	tParent.Run("buy and sell on assets", func(t *testing.T) {
		t.Parallel()

		openOrders := []models.Order{
			models.Order{
				OrderID: uuid1,
				Side: "buy",
				Amount: 200,
				ChildOrders: []models.Order{
					models.Order{ Amount: 100, Side: "buy" },
				},
			},
			models.Order{
				OrderID: uuid2,
				Side: "sell",
				Amount: 200,
				ChildOrders: []models.Order{
					models.Order{ Amount: 100, Side: "sell" },
				},
			},
		}

		result := srvc.GetOrders(openOrders, 100, 100)

		assert.Equal(t, len(result), 2)

		assert.Equal(t, result[0].Side, "buy")
		assert.Equal(t, result[0].Amount, 100)
		assert.Equal(t, *result[0].ParentOrderID, uuid1)

		assert.Equal(t, result[1].Side, "sell")
		assert.Equal(t, result[1].Amount, 100)
		assert.Equal(t, *result[1].ParentOrderID, uuid2)
	})

	tParent.Run("buy and sell on match", func(t *testing.T) {
		t.Parallel()

		openOrders := []models.Order{
			models.Order{
				OrderID: uuid1,
				Side: "buy",
				Amount: 100,
				ChildOrders: []models.Order{},
			},
			models.Order{
				OrderID: uuid2,
				Side: "sell",
				Amount: 100,
				ChildOrders: []models.Order{},
			},
		}

		result := srvc.GetOrders(openOrders, 0, 0)

		assert.Equal(t, len(result), 2)

		assert.Equal(t, result[0].Side, "buy")
		assert.Equal(t, result[0].Amount, 100)
		assert.Equal(t, *result[0].ParentOrderID, uuid1)
		assert.Equal(t, *result[0].MatchingOrderID, result[1].OrderID)

		assert.Equal(t, result[1].Side, "sell")
		assert.Equal(t, result[1].Amount, 100)
		assert.Equal(t, *result[1].ParentOrderID, uuid2)
		assert.Equal(t, *result[1].MatchingOrderID, result[0].OrderID)
	})
}
