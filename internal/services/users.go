package services

import (
	"go_server/internal/db"
	"go_server/internal/models"

	"github.com/google/uuid"
)

func GetUser(userID uuid.UUID) models.User {
  var user models.User

  err := db.DB.First(&user, userID).Error

  if err != nil {
    panic("Error finding user")
  }

	return user
}

func ListUsers(query *models.User) []models.User {
  var users []models.User

  err := db.DB.Where(query).Order("created_at desc").Find(&users).Error

  if err != nil {
    panic("Error listing users")
  }

	return users
}

func CreateUser(userPayload models.User) models.User {
	user := userPayload

  err := db.DB.Create(&user).Error

  if err != nil {
    panic("Error creating user")
  }

	return user
}

func ModifyUser(userID uuid.UUID, userPayload models.User) models.User {
	var userFound models.User

	errFind := db.DB.Where("user_id = ?", userID).First(&userFound).Error

	if errFind != nil {
		panic("Error finding user")
	}

	if (userPayload.FirstName != nil) {
		userFound.FirstName = userPayload.FirstName
	}

	if (userPayload.LastName != nil) {
		userFound.LastName = userPayload.LastName
	}

	if (userPayload.Auth0ID != nil) {
		userFound.Auth0ID = userPayload.Auth0ID
	}

  errUpdate := db.DB.Save(&userFound).Error

  if errUpdate != nil {
    panic("Error updating user")
  }

	return userFound
}

func DeleteUser(userID uuid.UUID) {

  errUpdate := db.DB.Delete(&models.User{}, userID).Error

  if errUpdate != nil {
    panic("Error deleting user")
  }
}
