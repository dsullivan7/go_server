package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetPortfolioIndustry(portfolioIndustryID uuid.UUID) (*models.PortfolioIndustry, error) {
	var portfolioIndustry models.PortfolioIndustry

	err := gormStore.database.First(&portfolioIndustry, portfolioIndustryID).Error
	if err != nil {
		return nil, err
	}

	return &portfolioIndustry, nil
}

func (gormStore *Store) ListPortfolioIndustries(query map[string]interface{}) ([]models.PortfolioIndustry, error) {
	var portfolioIndustries []models.PortfolioIndustry

	err := gormStore.database.Where(query).Order("created_at desc").Find(&portfolioIndustries).Error
	if err != nil {
		return nil, err
	}

	return portfolioIndustries, nil
}

func (gormStore *Store) CreatePortfolioIndustry(
	portfolioIndustryPayload models.PortfolioIndustry,
) (*models.PortfolioIndustry, error) {
	portfolioIndustry := portfolioIndustryPayload

	err := gormStore.database.Create(&portfolioIndustry).Error
	if err != nil {
		return nil, err
	}

	return &portfolioIndustry, nil
}

func (gormStore *Store) ModifyPortfolioIndustry(
	portfolioIndustryID uuid.UUID,
	portfolioIndustryPayload models.PortfolioIndustry,
) (*models.PortfolioIndustry, error) {
	var portfolioIndustryFound models.PortfolioIndustry

	errFind := gormStore.database.Where(
		"portfolio_industry_id = ?",
		portfolioIndustryID,
	).First(&portfolioIndustryFound).Error

	if errFind != nil {
		return nil, errFind
	}

	errUpdate := gormStore.database.Save(&portfolioIndustryFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &portfolioIndustryFound, nil
}

func (gormStore *Store) DeletePortfolioIndustry(portfolioIndustryID uuid.UUID) error {
	err := gormStore.database.Delete(&models.PortfolioIndustry{}, portfolioIndustryID).Error

	return err
}
