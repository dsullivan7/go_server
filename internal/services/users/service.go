package users

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func Get(userID uuid.UUID) models.User {
	user := models.User{
		UserID: userID,
	}

	return user
}
