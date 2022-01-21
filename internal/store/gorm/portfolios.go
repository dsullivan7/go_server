package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetPortfolio(portfolioID uuid.UUID) (*models.Portfolio, error) {
	var portfolio models.Portfolio

	err := gormStore.database.First(&portfolio, portfolioID).Error
	if err != nil {
		return nil, err
	}

	return &portfolio, nil
}

func (gormStore *Store) ListPortfolios(query map[string]interface{}) ([]models.Portfolio, error) {
	var portfolios []models.Portfolio

	err := gormStore.database.Where(query).Order("created_at desc").Find(&portfolios).Error
	if err != nil {
		return nil, err
	}

	return portfolios, nil
}

func (gormStore *Store) CreatePortfolio(portfolioPayload models.Portfolio) (*models.Portfolio, error) {
	portfolio := portfolioPayload

	err := gormStore.database.Create(&portfolio).Error
	if err != nil {
		return nil, err
	}

	return &portfolio, nil
}

func (gormStore *Store) ModifyPortfolio(
	portfolioID uuid.UUID,
	portfolioPayload models.Portfolio,
) (*models.Portfolio, error) {
	var portfolioFound models.Portfolio

	errFind := gormStore.database.Where("portfolio_id = ?", portfolioID).First(&portfolioFound).Error

	if errFind != nil {
		return nil, errFind
	}

	if portfolioPayload.Risk != 0 {
		portfolioFound.Risk = portfolioPayload.Risk
	}

	errUpdate := gormStore.database.Save(&portfolioFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &portfolioFound, nil
}

func (gormStore *Store) DeletePortfolio(portfolioID uuid.UUID) error {
	err := gormStore.database.Delete(&models.Portfolio{}, portfolioID).Error

	return err
}
