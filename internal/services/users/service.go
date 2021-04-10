package users

import (
  "github.com/satori/go.uuid"

  "go_server/internal/models"
)

func Get(userID uuid.UUID) models.User{
  user := models.User {
    UserID: userID,
  }

  return user
}
