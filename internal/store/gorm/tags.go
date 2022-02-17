package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetTag(tagID uuid.UUID) (*models.Tag, error) {
	var tag models.Tag

	err := gormStore.database.First(&tag, tagID).Error
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (gormStore *Store) ListTags(query map[string]interface{}) ([]models.Tag, error) {
	var tags []models.Tag

	err := gormStore.database.Where(query).Order("created_at desc").Find(&tags).Error
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (gormStore *Store) CreateTag(tagPayload models.Tag) (*models.Tag, error) {
	tag := tagPayload

	err := gormStore.database.Create(&tag).Error
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (gormStore *Store) ModifyTag(
	tagID uuid.UUID,
	tagPayload models.Tag,
) (*models.Tag, error) {
	var tagFound models.Tag

	errFind := gormStore.database.Where("tag_id = ?", tagID).First(&tagFound).Error

	if errFind != nil {
		return nil, errFind
	}

	if tagPayload.Name != nil {
		tagFound.Name = tagPayload.Name
	}

	errUpdate := gormStore.database.Save(&tagFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &tagFound, nil
}

func (gormStore *Store) DeleteTag(tagID uuid.UUID) error {
	err := gormStore.database.Delete(&models.Tag{}, tagID).Error

	return err
}
