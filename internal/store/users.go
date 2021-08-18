package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *GormStore) GetUser(userID uuid.UUID) models.User {
  var user models.User

  err := gormStore.database.First(&user, userID).Error

  if err != nil {
    panic("Error finding user")
  }

	return user
}

func (gormStore *GormStore) ListUsers(query map[string]interface{}) []models.User {
  var users []models.User

  err := gormStore.database.Where(query).Order("created_at desc").Find(&users).Error

  if err != nil {
    panic("Error listing users")
  }

	return users
}

func (gormStore *GormStore) CreateUser(userPayload models.User) models.User {
	user := userPayload

  err := gormStore.database.Create(&user).Error

  if err != nil {
    panic("Error creating user")
  }

	return user
}

func (gormStore *GormStore) ModifyUser(userID uuid.UUID, userPayload models.User) models.User {
	var userFound models.User

	errFind := gormStore.database.Where("user_id = ?", userID).First(&userFound).Error

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

  errUpdate := gormStore.database.Save(&userFound).Error

  if errUpdate != nil {
    panic("Error updating user")
  }

	return userFound
}

func (gormStore *GormStore) DeleteUser(userID uuid.UUID) {

  errUpdate := gormStore.database.Delete(&models.User{}, userID).Error

  if errUpdate != nil {
    panic("Error deleting user")
  }
}
