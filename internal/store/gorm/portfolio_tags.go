package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetPortfolioTag(portfolioTagID uuid.UUID) (*models.PortfolioTag, error) {
	var portfolioTag models.PortfolioTag

	err := gormStore.database.First(&portfolioTag, portfolioTagID).Error
	if err != nil {
		return nil, err
	}

	return &portfolioTag, nil
}

func (gormStore *Store) ListPortfolioTags(query map[string]interface{}) ([]models.PortfolioTag, error) {
	var portfolioTags []models.PortfolioTag

	err := gormStore.database.Where(query).Order("created_at desc").Find(&portfolioTags).Error
	if err != nil {
		return nil, err
	}

	return portfolioTags, nil
}

func (gormStore *Store) CreatePortfolioTag(
	portfolioTagPayload models.PortfolioTag,
) (*models.PortfolioTag, error) {
	portfolioTag := portfolioTagPayload

	err := gormStore.database.Create(&portfolioTag).Error
	if err != nil {
		return nil, err
	}

	return &portfolioTag, nil
}

func (gormStore *Store) ModifyPortfolioTag(
	portfolioTagID uuid.UUID,
	portfolioTagPayload models.PortfolioTag,
) (*models.PortfolioTag, error) {
	var portfolioTagFound models.PortfolioTag

	errFind := gormStore.database.Where(
		"portfolio_tag_id = ?",
		portfolioTagID,
	).First(&portfolioTagFound).Error

	if errFind != nil {
		return nil, errFind
	}

	errUpdate := gormStore.database.Save(&portfolioTagFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &portfolioTagFound, nil
}

func (gormStore *Store) DeletePortfolioTag(portfolioTagID uuid.UUID) error {
	err := gormStore.database.Delete(&models.PortfolioTag{}, portfolioTagID).Error

	return err
}
