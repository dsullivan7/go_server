package services

import (
	"go_server/internal/models"
	"time"
)

func (mockService *MockService) GetOrders(
	openOrders []models.Order,
	childOrders []models.Order,
	netSecurityValue int,
	netCashValue int,
) []models.Order {
	args := mockService.Called(openOrders, childOrders, netSecurityValue, netCashValue)

	return args.Get(0).([]models.Order)
}

 func (mockService *MockService) GetBalance(orders []models.Order, interest float64, currentTime time.Time) (int, int, int) {
	 args := mockService.Called(orders, interest, currentTime)

	 return args.Int(0), args.Int(1), args.Int(2)
 }
