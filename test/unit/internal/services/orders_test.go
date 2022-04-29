package services_test

import (
	"time"

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
	uuid3 := uuid.New()

	tParent.Run("buy and sell on assets", func(t *testing.T) {
		t.Parallel()

		openOrders := []models.Order{
			models.Order{
				OrderID: uuid1,
				Side:    "buy",
				Amount:  200,
			},
			models.Order{
				OrderID: uuid2,
				Side:    "sell",
				Amount:  200,
			},
		}

		childOrders := []models.Order{
			models.Order{Amount: 100, Side: "buy", ParentOrderID: &uuid1},
			models.Order{Amount: 100, Side: "sell", ParentOrderID: &uuid2},
		}

		result := srvc.GetOrders(openOrders, childOrders, 100, 100)

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
				Side:    "buy",
				Amount:  100,
			},
			models.Order{
				OrderID: uuid2,
				Side:    "sell",
				Amount:  100,
			},
		}

		result := srvc.GetOrders(openOrders, []models.Order{}, 0, 0)

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

	tParent.Run("mixed buy and sell on match and assets", func(t *testing.T) {
		t.Parallel()

		openOrders := []models.Order{
			models.Order{
				OrderID: uuid1,
				Side:    "buy",
				Amount:  200,
			},
			models.Order{
				OrderID: uuid2,
				Side:    "sell",
				Amount:  200,
			},
			models.Order{
				OrderID: uuid3,
				Side:    "sell",
				Amount:  200,
			},
		}

		childOrders := []models.Order{
			models.Order{
				ParentOrderID: &uuid1,
				Side:          "buy",
				Amount:        50,
			},
			models.Order{
				ParentOrderID: &uuid2,
				Side:          "sell",
				Amount:        150,
			},
		}

		result := srvc.GetOrders(openOrders, childOrders, 100, 100)

		assert.Equal(t, len(result), 5)

		assert.Equal(t, result[0].Side, "buy")
		assert.Equal(t, result[0].Amount, 100)
		assert.Equal(t, *result[0].ParentOrderID, uuid1)

		assert.Equal(t, result[1].Side, "sell")
		assert.Equal(t, result[1].Amount, 50)
		assert.Equal(t, *result[1].ParentOrderID, uuid2)

		assert.Equal(t, result[2].Side, "sell")
		assert.Equal(t, result[2].Amount, 50)
		assert.Equal(t, *result[2].ParentOrderID, uuid3)

		assert.Equal(t, result[3].Side, "buy")
		assert.Equal(t, result[3].Amount, 50)
		assert.Equal(t, *result[3].ParentOrderID, uuid1)
		assert.Equal(t, *result[3].MatchingOrderID, result[4].OrderID)

		assert.Equal(t, result[4].Side, "sell")
		assert.Equal(t, result[4].Amount, 50)
		assert.Equal(t, *result[4].ParentOrderID, uuid3)
		assert.Equal(t, *result[4].MatchingOrderID, result[3].OrderID)
	})

	tParent.Run("GetReturn simple interest", func(t *testing.T) {
		t.Parallel()

		time1 := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
		time2 := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)

		orders := []models.Order{
			models.Order{
				OrderID:     uuid1,
				Side:        "buy",
				Amount:      100,
				CompletedAt: &time1,
			},
		}

		total, principal, interest := srvc.GetBalance(orders, 0.05, time2)

		assert.Equal(t, total, 105)
		assert.Equal(t, principal, 100)
		assert.Equal(t, interest, 5)
	})

	tParent.Run("GetReturn simple interest buy and sell", func(t *testing.T) {
		t.Parallel()

		time1 := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
		time2 := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
		time3 := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

		orders := []models.Order{
			models.Order{
				OrderID:     uuid1,
				Side:        "buy",
				Amount:      10000,
				CompletedAt: &time1,
			},
			models.Order{
				OrderID:     uuid1,
				Side:        "sell",
				Amount:      5000,
				CompletedAt: &time2,
			},
		}

		total, principal, interest := srvc.GetBalance(orders, 0.05, time3)

		assert.Equal(t, total, 5750)
		assert.Equal(t, principal, 5000)
		assert.Equal(t, interest, 750)
	})
}
