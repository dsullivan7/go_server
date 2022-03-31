package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetItem(itemID uuid.UUID) (*models.Item, error) {
	var item models.Item

	err := gormStore.database.First(&item, itemID).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (gormStore *Store) ListItems(query map[string]interface{}) ([]models.Item, error) {
	var items []models.Item

	err := gormStore.database.Where(query).Order("created_at desc").Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (gormStore *Store) CreateItem(itemPayload models.Item) (*models.Item, error) {
	item := itemPayload

	err := gormStore.database.Create(&item).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (gormStore *Store) ModifyItem(
	itemID uuid.UUID,
	itemPayload models.Item,
) (*models.Item, error) {
	var itemFound models.Item

	errFind := gormStore.database.Where("item_id = ?", itemID).First(&itemFound).Error

	if errFind != nil {
		return nil, errFind
	}

	if itemPayload.Name != "" {
		itemFound.Name = itemPayload.Name
	}

	errUpdate := gormStore.database.Save(&itemFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &itemFound, nil
}

func (gormStore *Store) DeleteItem(itemID uuid.UUID) error {
	err := gormStore.database.Delete(&models.Item{}, itemID).Error

	return err
}
