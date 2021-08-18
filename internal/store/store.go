package store

import (
  "go_server/internal/models"

  "gorm.io/gorm"
  "github.com/google/uuid"
)

type Store interface {
  GetUser(userID uuid.UUID) models.User
  ListUsers(query map[string]interface{}) []models.User
  CreateUser(userPayload models.User) models.User
  ModifyUser(userID uuid.UUID, userPayload models.User) models.User
  DeleteUser(userID uuid.UUID)

  GetReview(userID uuid.UUID) models.Review
  ListReviews(query map[string]interface{}) []models.Review
  CreateReview(userPayload models.Review) models.Review
  ModifyReview(userID uuid.UUID, userPayload models.Review) models.Review
  DeleteReview(userID uuid.UUID)
}

type GormStore struct {
	database *gorm.DB
}

func NewGormStore(database *gorm.DB) Store {
  return &GormStore{ database: database }
}
