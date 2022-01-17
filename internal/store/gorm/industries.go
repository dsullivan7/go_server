package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetIndustry(industryID uuid.UUID) (*models.Industry, error) {
	var industry models.Industry

	err := gormStore.database.First(&industry, industryID).Error
	if err != nil {
		return nil, err
	}

	return &industry, nil
}

func (gormStore *Store) ListIndustries(query map[string]interface{}) ([]models.Industry, error) {
	var industries []models.Industry

	err := gormStore.database.Where(query).Order("created_at desc").Find(&industries).Error
	if err != nil {
		return nil, err
	}

	return industries, nil
}

func (gormStore *Store) CreateIndustry(industryPayload models.Industry) (*models.Industry, error) {
	industry := industryPayload

	err := gormStore.database.Create(&industry).Error
	if err != nil {
		return nil, err
	}

	return &industry, nil
}

func (gormStore *Store) ModifyIndustry(
	industryID uuid.UUID,
	industryPayload models.Industry,
) (*models.Industry, error) {
	var industryFound models.Industry

	errFind := gormStore.database.Where("industry_id = ?", industryID).First(&industryFound).Error

	if errFind != nil {
		return nil, errFind
	}

	if industryPayload.Name != nil {
		industryFound.Name = industryPayload.Name
	}

	errUpdate := gormStore.database.Save(&industryFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &industryFound, nil
}

func (gormStore *Store) DeleteIndustry(industryID uuid.UUID) error {
	err := gormStore.database.Delete(&models.Industry{}, industryID).Error

	return err
}
