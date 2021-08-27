package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetReview(reviewID uuid.UUID) (*models.Review, error) {
	var review models.Review

	err := gormStore.database.First(&review, reviewID).Error
	if err != nil {
		return nil, err
	}

	return &review, nil
}

func (gormStore *Store) ListReviews(query map[string]interface{}) ([]models.Review, error) {
	var reviews []models.Review

	err := gormStore.database.Where(query).Order("created_at desc").Find(&reviews).Error
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (gormStore *Store) CreateReview(reviewPayload models.Review) (*models.Review, error) {
	review := reviewPayload

	err := gormStore.database.Create(&review).Error
	if err != nil {
		return nil, err
	}

	return &review, nil
}

func (gormStore *Store) ModifyReview(reviewID uuid.UUID, reviewPayload models.Review) (*models.Review, error) {
	var reviewFound models.Review

	errFind := gormStore.database.Where("review_id = ?", reviewID).First(&reviewFound).Error

	if errFind != nil {
		return nil, errFind
	}

	if reviewPayload.FromUserID != nil {
		reviewFound.FromUserID = reviewPayload.FromUserID
	}

	if reviewPayload.ToUserID != nil {
		reviewFound.ToUserID = reviewPayload.ToUserID
	}

	if reviewPayload.Text != nil {
		reviewFound.Text = reviewPayload.Text
	}

	errUpdate := gormStore.database.Save(&reviewFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &reviewFound, nil
}

func (gormStore *Store) DeleteReview(reviewID uuid.UUID) error {
	err := gormStore.database.Delete(&models.Review{}, reviewID).Error

	return err
}
