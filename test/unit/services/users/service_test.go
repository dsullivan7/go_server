package users_test

import (
	"testing"

	"go_server/internal/db"
	"go_server/internal/models"
	UsersService "go_server/internal/services/users"
)

func TestGet(t *testing.T) {
	t.Parallel()

  userCreated := models.User{ FirstName: "FirstName" }
  db.Connect()
  db.DB.Create(&userCreated)

	userFound := UsersService.Get(userCreated.UserID)

	if userFound.UserID != userCreated.UserID {
		t.Errorf("userID incorrect, got: %s, want: %s.", userFound.UserID, userCreated.UserID)
	}
}
