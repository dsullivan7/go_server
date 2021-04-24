package users

import (
	"go_server/internal/db"
	"go_server/internal/models"

	"github.com/google/uuid"
)

func Get(userID uuid.UUID) models.User {
  var user models.User

  err := db.DB.Where("user_id = ?", userID).First(&user).Error

  if err != nil {
    panic("Error finding user")
  }

	return user
}

func List() []models.User {
  var users []models.User

  err := db.DB.Find(&users).Error

  if err != nil {
    panic("Error listing users")
  }

	return users
}

func Create(userPayload models.User) models.User {
	user := userPayload

  err := db.DB.Create(&user).Error

  if err != nil {
    panic("Error creating user")
  }

	return user
}
