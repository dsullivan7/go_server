package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetUserIndustry(userIndustryID uuid.UUID) (*models.UserIndustry, error) {
	var userIndustry models.UserIndustry

	err := gormStore.database.First(&userIndustry, userIndustryID).Error
	if err != nil {
		return nil, err
	}

	return &userIndustry, nil
}

func (gormStore *Store) ListUserIndustries(query map[string]interface{}) ([]models.UserIndustry, error) {
	var userIndustries []models.UserIndustry

	err := gormStore.database.Where(query).Order("created_at desc").Find(&userIndustries).Error
	if err != nil {
		return nil, err
	}

	return userIndustries, nil
}

func (gormStore *Store) CreateUserIndustry(userIndustryPayload models.UserIndustry) (*models.UserIndustry, error) {
	userIndustry := userIndustryPayload

	err := gormStore.database.Create(&userIndustry).Error
	if err != nil {
		return nil, err
	}

	return &userIndustry, nil
}

func (gormStore *Store) ModifyUserIndustry(
	userIndustryID uuid.UUID,
	userIndustryPayload models.UserIndustry,
) (*models.UserIndustry, error) {
	var userIndustryFound models.UserIndustry

	errFind := gormStore.database.Where("user_industry_id = ?", userIndustryID).First(&userIndustryFound).Error

	if errFind != nil {
		return nil, errFind
	}

	errUpdate := gormStore.database.Save(&userIndustryFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &userIndustryFound, nil
}

func (gormStore *Store) DeleteUserIndustry(userIndustryID uuid.UUID) error {
	err := gormStore.database.Delete(&models.UserIndustry{}, userIndustryID).Error

	return err
}
