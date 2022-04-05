package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetGroup(groupID uuid.UUID) (*models.Group, error) {
	var group models.Group

	err := gormStore.database.First(&group, groupID).Error
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (gormStore *Store) ListGroups(query map[string]interface{}) ([]models.Group, error) {
	var groups []models.Group

	err := gormStore.database.Where(query).Order("created_at desc").Find(&groups).Error
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (gormStore *Store) CreateGroup(groupPayload models.Group) (*models.Group, error) {
	group := groupPayload

	err := gormStore.database.Create(&group).Error
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (gormStore *Store) ModifyGroup(
	groupID uuid.UUID,
	groupPayload models.Group,
) (*models.Group, error) {
	var groupFound models.Group

	errFind := gormStore.database.Where("group_id = ?", groupID).First(&groupFound).Error

	if errFind != nil {
		return nil, errFind
	}

	if groupPayload.Name != "" {
		groupFound.Name = groupPayload.Name
	}

	errUpdate := gormStore.database.Save(&groupFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &groupFound, nil
}

func (gormStore *Store) DeleteGroup(groupID uuid.UUID) error {
	err := gormStore.database.Delete(&models.Group{}, groupID).Error

	return err
}
