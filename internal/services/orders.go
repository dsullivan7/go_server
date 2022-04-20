package services

import (
	"math"
	"go_server/internal/models"
)

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

	// if netSecurityValue is greater than zero, create a child order to take the remaining asset value
	if (netSecurityValue > 0) {
		remainingSecurityValue := netSecurityValue
		for _, openOrderBuy := range openOrdersBuy {
			// break the loop if no security value remains
			if (remainingSecurityValue <= 0) {
				break
			}

			// calculate the total amount of the order that has already been covered
			totalOrderCovered := 0
			for _, childOrder := range openOrderBuy.ChildOrders {
				totalOrderCovered += childOrder.Amount
			}

			// set the amount of the child order to be the min of the remaining security value and the amount left in the parent order 
			childOrderAmount := int(math.Min(float64(remainingSecurityValue), float64(openOrderBuy.Amount - totalOrderCovered)))

			childOrder := models.Order{
				ParentOrderID: &openOrderBuy.OrderID,
				Amount: childOrderAmount,
				Side: "buy",
			}

			openOrderBuy.ChildOrders = append(openOrderBuy.ChildOrders, childOrder)
			orders = append(orders, childOrder)
			remainingSecurityValue -= childOrderAmount
		}
	}

	return orders
}
