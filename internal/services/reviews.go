package services

import (
	"go_server/internal/db"
	"go_server/internal/models"

	"github.com/google/uuid"
)

func GetReview(reviewID uuid.UUID) models.Review {
  var review models.Review

  err := db.DB.First(&review, reviewID).Error

  if err != nil {
    panic("Error finding review")
  }

	return review
}

func ListReviews(query *models.Review) []models.Review {
  var reviews []models.Review

  err := db.DB.Where(query).Order("created_at desc").Find(&reviews).Error

  if err != nil {
    panic("Error listing reviews")
  }

	return reviews
}

func CreateReview(reviewPayload models.Review) models.Review {
	review := reviewPayload

  err := db.DB.Create(&review).Error

  if err != nil {
    panic("Error creating review")
  }

	return review
}

func ModifyReview(reviewID uuid.UUID, reviewPayload models.Review) models.Review {
	var reviewFound models.Review

	errFind := db.DB.Where("review_id = ?", reviewID).First(&reviewFound).Error

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

  errUpdate := db.DB.Save(&reviewFound).Error

  if errUpdate != nil {
    panic("Error updating review")
  }

	return reviewFound
}

func DeleteReview(reviewID uuid.UUID) {

  errUpdate := db.DB.Delete(&models.Review{}, reviewID).Error

  if errUpdate != nil {
    panic("Error deleting review")
  }
}
