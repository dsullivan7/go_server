package services

import (
	"math"
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (srvc *Service) getAmountRemaining(openOrder models.Order) int {
	orderRemaining := openOrder.Amount

	for _, childOrder := range openOrder.ChildOrders {
		orderRemaining -= childOrder.Amount
	}

	return orderRemaining
}

func (srvc *Service) getAssetOrders(openOrders []models.Order, val int, side string) []models.Order {
	var orders []models.Order
	// if the target val is greater than zero, create a child order to take the remaining asset value
	if (val > 0) {
		remainingSecurityValue := val
		for _, openOrder := range openOrders {
			// break the loop if no security value remains
			if (remainingSecurityValue <= 0) {
				break
			}

			// set the amount of the child order to be
			// the min of the remaining security value and the amount left in the parent order
			childOrderAmount := int(math.Min(float64(remainingSecurityValue), float64(srvc.getAmountRemaining(openOrder))))

			parentOrderID := openOrder.OrderID

			childOrder := models.Order{
				ParentOrderID: &parentOrderID,
				Amount: childOrderAmount,
				Side: side,
			}

			// append the new child order to the open order
			orders = append(orders, childOrder)
			remainingSecurityValue -= childOrderAmount
		}
	}

	return orders
}

func (srvc *Service) getMatchOrders(openBuyOrders []models.Order, openSellOrders []models.Order) []models.Order {
	var orders []models.Order

	i := 0
	j := 0

	for (i < len(openBuyOrders) && j < len(openSellOrders)) {
		openBuyOrder := &openBuyOrders[i]
		openSellOrder := &openSellOrders[j]

		remainingBuy := srvc.getAmountRemaining(*openBuyOrder)
		remainingSell := srvc.getAmountRemaining(*openSellOrder)

		if remainingBuy == 0 {
			i++

			continue
		}

		if (remainingSell == 0) {
			j++

			continue
		}

		uuid1 := uuid.New()
		uuid2 := uuid.New()

		childOrderAmount := int(math.Min(float64(remainingBuy), float64(remainingSell)))

		childOrder1 := models.Order{
			OrderID: uuid1,
			ParentOrderID: &openBuyOrder.OrderID,
			MatchingOrderID: &uuid2,
			Amount: childOrderAmount,
			Side: "buy",
		}

		childOrder2 := models.Order{
			OrderID: uuid2,
			ParentOrderID: &openSellOrder.OrderID,
			MatchingOrderID: &uuid1,
			Amount: childOrderAmount,
			Side: "sell",
		}

		openBuyOrder.ChildOrders = append(openBuyOrder.ChildOrders, childOrder1)
		openSellOrder.ChildOrders = append(openBuyOrder.ChildOrders, childOrder2)
		orders = append(orders, childOrder1, childOrder2)
	}

	return orders
}

// ListOrders returns what orders need to be made to resolve the market under the given parameters.
func (srvc *Service) GetOrders(
  openOrders []models.Order,
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

	childBuyOrders := srvc.getAssetOrders(openOrdersBuy, netSecurityValue, "buy")
	orders = append(orders, childBuyOrders...)

	for _, childOrder := range childBuyOrders {
		for i := range openOrdersBuy {
			if *childOrder.ParentOrderID == openOrdersBuy[i].OrderID {
				openOrdersBuy[i].ChildOrders = append(openOrdersBuy[i].ChildOrders, childOrder)

				break
			}
		}
	}

	childSellOrders := srvc.getAssetOrders(openOrdersSell, netCashValue, "sell")
	orders = append(orders, childSellOrders...)

	for _, childOrder := range childSellOrders {
		for i := range openOrdersSell {
			if *childOrder.ParentOrderID == openOrdersSell[i].OrderID {
				openOrdersSell[i].ChildOrders = append(openOrdersSell[i].ChildOrders, childOrder)

				break
			}
		}
	}

	childMatchOrders := srvc.getMatchOrders(openOrdersBuy, openOrdersSell)
	orders = append(orders, childMatchOrders...)

	return orders
}
