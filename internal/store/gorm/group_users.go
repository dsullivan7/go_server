package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetGroupUser(groupUserID uuid.UUID) (*models.GroupUser, error) {
	var groupUser models.GroupUser

	err := gormStore.database.First(&groupUser, groupUserID).Error
	if err != nil {
		return nil, err
	}

	return &groupUser, nil
}

func (gormStore *Store) ListGroupUsers(query map[string]interface{}) ([]models.GroupUser, error) {
	var groupUsers []models.GroupUser

	err := gormStore.database.Where(query).Order("created_at desc").Find(&groupUsers).Error
	if err != nil {
		return nil, err
	}

	return groupUsers, nil
}

func (gormStore *Store) CreateGroupUser(groupUserPayload models.GroupUser) (*models.GroupUser, error) {
	groupUser := groupUserPayload

	err := gormStore.database.Create(&groupUser).Error
	if err != nil {
		return nil, err
	}

	return &groupUser, nil
}

func (gormStore *Store) ModifyGroupUser(
	groupUserID uuid.UUID,
	groupUserPayload models.GroupUser,
) (*models.GroupUser, error) {
	var groupUserFound models.GroupUser

	errFind := gormStore.database.Where("group_user_id = ?", groupUserID).First(&groupUserFound).Error

	if errFind != nil {
		return nil, errFind
	}

	errUpdate := gormStore.database.Save(&groupUserFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &groupUserFound, nil
}

func (gormStore *Store) DeleteGroupUser(groupUserID uuid.UUID) error {
	err := gormStore.database.Delete(&models.GroupUser{}, groupUserID).Error

	return err
}
