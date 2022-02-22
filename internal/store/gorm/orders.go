package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetOrder(orderID uuid.UUID) (*models.Order, error) {
	var order models.Order

	err := gormStore.database.First(&order, orderID).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (gormStore *Store) ListOrders(query map[string]interface{}) ([]models.Order, error) {
	var orders []models.Order

	err := gormStore.database.Where(query).Order("created_at desc").Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (gormStore *Store) CreateOrder(orderPayload models.Order) (*models.Order, error) {
	order := orderPayload

	err := gormStore.database.Create(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (gormStore *Store) ModifyOrder(
	orderID uuid.UUID,
	orderPayload models.Order,
) (*models.Order, error) {
	var orderFound models.Order

	errFind := gormStore.database.Where("order_id = ?", orderID).First(&orderFound).Error

	if errFind != nil {
		return nil, errFind
	}

	errUpdate := gormStore.database.Save(&orderFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &orderFound, nil
}

func (gormStore *Store) DeleteOrder(orderID uuid.UUID) error {
	err := gormStore.database.Delete(&models.Order{}, orderID).Error

	return err
}
