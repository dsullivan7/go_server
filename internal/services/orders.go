package services

import (
	"math"
	"time"

	"go_server/internal/models"

	"github.com/google/uuid"
)

const daysInYear = 365
const hoursInDay = 24

func (srvc *Service) getAmountRemaining(openOrder models.Order, childOrders []models.Order) int {
	orderRemaining := openOrder.Amount

	for _, childOrder := range childOrders {
		if *childOrder.ParentOrderID == openOrder.OrderID {
			orderRemaining -= childOrder.Amount
		}
	}

	return orderRemaining
}

func (srvc *Service) getAssetOrders(
	openOrders []models.Order,
	childOrders []models.Order,
	val int,
	side string,
) []models.Order {
	var orders []models.Order

	// if the target val is greater than zero, create a child order to take the remaining asset value
	if val > 0 {
		remainingSecurityValue := val
		for _, openOrder := range openOrders {
			// break the loop if no security value remains
			if remainingSecurityValue <= 0 {
				break
			}

			// set the amount of the child order to be
			// the min of the remaining security value and the amount left in the parent order
			childOrderAmount := int(math.Min(
				float64(remainingSecurityValue),
				float64(srvc.getAmountRemaining(openOrder, childOrders)),
			))

			parentOrderID := openOrder.OrderID

			childOrder := models.Order{
				OrderID:         uuid.New(),
				MatchingOrderID: nil,
				ParentOrderID:   &parentOrderID,
				Amount:          childOrderAmount,
				Side:            side,
			}

			// append the new child order to the open order
			orders = append(orders, childOrder)
			remainingSecurityValue -= childOrderAmount
		}
	}

	return orders
}

func (srvc *Service) getMatchOrders(
	openBuyOrders []models.Order,
	openSellOrders []models.Order,
	childOrders []models.Order,
) []models.Order {
	var orders []models.Order

	i := 0
	j := 0

	for i < len(openBuyOrders) && j < len(openSellOrders) {
		openBuyOrder := openBuyOrders[i]
		openSellOrder := openSellOrders[j]

		remainingBuy := srvc.getAmountRemaining(openBuyOrder, childOrders)
		remainingSell := srvc.getAmountRemaining(openSellOrder, childOrders)

		if remainingBuy == 0 {
			i++

			continue
		}

		if remainingSell == 0 {
			j++

			continue
		}

		uuid1 := uuid.New()
		uuid2 := uuid.New()

		childOrderAmount := int(math.Min(float64(remainingBuy), float64(remainingSell)))

		childOrder1 := models.Order{
			OrderID:         uuid1,
			ParentOrderID:   &openBuyOrder.OrderID,
			MatchingOrderID: &uuid2,
			Amount:          childOrderAmount,
			Side:            "buy",
		}

		childOrder2 := models.Order{
			OrderID:         uuid2,
			ParentOrderID:   &openSellOrder.OrderID,
			MatchingOrderID: &uuid1,
			Amount:          childOrderAmount,
			Side:            "sell",
		}

		orders = append(orders, childOrder1, childOrder2)
		childOrders = append(childOrders, childOrder1, childOrder2)
	}

	return orders
}

// ListOrders returns what orders need to be made to resolve the market under the given parameters.
func (srvc *Service) GetOrders(
	openOrders []models.Order,
	childOrders []models.Order,
	netSecurityValue int,
	netCashValue int,
) []models.Order {
	var orders []models.Order

	var openOrdersBuy []models.Order

	var openOrdersSell []models.Order

	for _, openOrder := range openOrders {
		switch openOrder.Side {
		case "buy":
			openOrdersBuy = append(openOrdersBuy, openOrder)
		case "sell":
			openOrdersSell = append(openOrdersSell, openOrder)
		}
	}

	childBuyOrders := srvc.getAssetOrders(openOrdersBuy, childOrders, netSecurityValue, "buy")
	orders = append(orders, childBuyOrders...)
	childOrders = append(childOrders, childBuyOrders...)

	childSellOrders := srvc.getAssetOrders(openOrdersSell, childOrders, netCashValue, "sell")
	orders = append(orders, childSellOrders...)
	childOrders = append(childOrders, childSellOrders...)

	childMatchOrders := srvc.getMatchOrders(openOrdersBuy, openOrdersSell, childOrders)
	orders = append(orders, childMatchOrders...)

	return orders
}

// GetBalance returns the balance, principal, and interest due to a user based on the orders and annual interest rate.
func (srvc *Service) GetBalance(orders []models.Order, interest float64, currentTime time.Time) (int, int, int) {
	var balance int

	var interestAmount float64

	for _, order := range orders {
		hours := currentTime.Sub(order.CompletedAt).Hours()
		orderInterest := float64(order.Amount) * hours * (interest / daysInYear / hoursInDay)

		switch order.Side {
		case "buy":
			balance += order.Amount
			interestAmount += orderInterest
		case "sell":
			balance -= order.Amount
			interestAmount -= orderInterest
		}
	}

	interestAmountRounded := int(math.Round(interestAmount))

	return (balance + interestAmountRounded), balance, interestAmountRounded
}
