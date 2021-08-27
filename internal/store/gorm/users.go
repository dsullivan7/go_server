package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetUser(userID uuid.UUID) (*models.User, error) {
	var user models.User

	err := gormStore.database.First(&user, userID).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (gormStore *Store) ListUsers(query map[string]interface{}) ([]models.User, error) {
	var users []models.User

	err := gormStore.database.Where(query).Order("created_at desc").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (gormStore *Store) CreateUser(userPayload models.User) (*models.User, error) {
	user := userPayload

	err := gormStore.database.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (gormStore *Store) ModifyUser(userID uuid.UUID, userPayload models.User) (*models.User, error) {
	var userFound models.User

	errFind := gormStore.database.Where("user_id = ?", userID).First(&userFound).Error

	if errFind != nil {
		return nil, errFind
	}

	if userPayload.FirstName != nil {
		userFound.FirstName = userPayload.FirstName
	}

	if userPayload.LastName != nil {
		userFound.LastName = userPayload.LastName
	}

	if userPayload.Auth0ID != nil {
		userFound.Auth0ID = userPayload.Auth0ID
	}

	err := gormStore.database.Save(&userFound).Error

	if err != nil {
		return nil, err
	}

	return &userFound, nil
}

func (gormStore *Store) DeleteUser(userID uuid.UUID) error {
	err := gormStore.database.Delete(&models.User{}, userID).Error

	return err
}
