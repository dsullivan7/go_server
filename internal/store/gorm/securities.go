package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetSecurity(securityID uuid.UUID) (*models.Security, error) {
	var security models.Security

	err := gormStore.database.First(&security, securityID).Error
	if err != nil {
		return nil, err
	}

	return &security, nil
}

func (gormStore *Store) ListSecurities(query map[string]interface{}) ([]models.Security, error) {
	var securities []models.Security

	err := gormStore.database.Where(query).Order("created_at desc").Find(&securities).Error
	if err != nil {
		return nil, err
	}

	return securities, nil
}

func (gormStore *Store) CreateSecurity(securityPayload models.Security) (*models.Security, error) {
	security := securityPayload

	err := gormStore.database.Create(&security).Error
	if err != nil {
		return nil, err
	}

	return &security, nil
}

func (gormStore *Store) ModifySecurity(
	securityID uuid.UUID,
	securityPayload models.Security,
) (*models.Security, error) {
	var securityFound models.Security

	errFind := gormStore.database.Where("security_id = ?", securityID).First(&securityFound).Error

	if errFind != nil {
		return nil, errFind
	}

	errUpdate := gormStore.database.Save(&securityFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &securityFound, nil
}

func (gormStore *Store) DeleteSecurity(securityID uuid.UUID) error {
	err := gormStore.database.Delete(&models.Security{}, securityID).Error

	return err
}
