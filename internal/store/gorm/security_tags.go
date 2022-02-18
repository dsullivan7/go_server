package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetSecurityTag(securityTagID uuid.UUID) (*models.SecurityTag, error) {
	var securityTag models.SecurityTag

	err := gormStore.database.First(&securityTag, securityTagID).Error
	if err != nil {
		return nil, err
	}

	return &securityTag, nil
}

func (gormStore *Store) ListSecurityTags(query map[string]interface{}) ([]models.SecurityTag, error) {
	var securityTags []models.SecurityTag

	err := gormStore.database.Where(query).Order("created_at desc").Find(&securityTags).Error
	if err != nil {
		return nil, err
	}

	return securityTags, nil
}

func (gormStore *Store) CreateSecurityTag(
	securityTagPayload models.SecurityTag,
) (*models.SecurityTag, error) {
	securityTag := securityTagPayload

	err := gormStore.database.Create(&securityTag).Error
	if err != nil {
		return nil, err
	}

	return &securityTag, nil
}

func (gormStore *Store) ModifySecurityTag(
	securityTagID uuid.UUID,
	securityTagPayload models.SecurityTag,
) (*models.SecurityTag, error) {
	var securityTagFound models.SecurityTag

	errFind := gormStore.database.Where(
		"security_tag_id = ?",
		securityTagID,
	).First(&securityTagFound).Error

	if errFind != nil {
		return nil, errFind
	}

	errUpdate := gormStore.database.Save(&securityTagFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &securityTagFound, nil
}

func (gormStore *Store) DeleteSecurityTag(securityTagID uuid.UUID) error {
	err := gormStore.database.Delete(&models.SecurityTag{}, securityTagID).Error

	return err
}
