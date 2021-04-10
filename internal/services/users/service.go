package users

import (
  "github.com/google/uuid"

  "go_server/internal/models"
)

func Get(userID uuid.UUID) models.User{
  user := models.User {
    UserID: userID,
  }

  return user
}
