package users

import (
  "go_server/internal/models"
)

func Get(userId string) models.User{
  user := models.User {
    UserId: userId,
  }

  return user
}
