package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *GormStore) GetReview(reviewID uuid.UUID) models.Review {
  var review models.Review

  err := gormStore.database.First(&review, reviewID).Error

  if err != nil {
    panic("Error finding review")
  }

	return review
}

func (gormStore *GormStore) ListReviews(query map[string]interface{}) []models.Review {
  var reviews []models.Review

  err := gormStore.database.Where(query).Order("created_at desc").Find(&reviews).Error

  if err != nil {
    panic("Error listing reviews")
  }

	return reviews
}

func (gormStore *GormStore) CreateReview(reviewPayload models.Review) models.Review {
	review := reviewPayload

  err := gormStore.database.Create(&review).Error

  if err != nil {
    panic("Error creating review")
  }

	return review
}

func (gormStore *GormStore) ModifyReview(reviewID uuid.UUID, reviewPayload models.Review) models.Review {
	var reviewFound models.Review

	errFind := gormStore.database.Where("review_id = ?", reviewID).First(&reviewFound).Error

	if errFind != nil {
		panic("Error finding review")
	}

	if (reviewPayload.FromUserID != nil) {
		reviewFound.FromUserID = reviewPayload.FromUserID
	}

	if (reviewPayload.ToUserID != nil) {
		reviewFound.ToUserID = reviewPayload.ToUserID
	}

	if (reviewPayload.Text != nil) {
		reviewFound.Text = reviewPayload.Text
	}

  errUpdate := gormStore.database.Save(&reviewFound).Error

  if errUpdate != nil {
    panic("Error updating review")
  }

	return reviewFound
}

func (gormStore *GormStore) DeleteReview(reviewID uuid.UUID) {

  errUpdate := gormStore.database.Delete(&models.Review{}, reviewID).Error

  if errUpdate != nil {
    panic("Error deleting review")
  }
}
