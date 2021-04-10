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
    panic("User not found")
  }

	return user
}
