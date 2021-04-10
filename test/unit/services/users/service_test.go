package users_test

import (
	"testing"

	"github.com/google/uuid"
	UsersService "go_server/internal/services/users"
)

func TestGet(t *testing.T) {
	userID := uuid.New()
	user := UsersService.Get(userID)
	if user.UserID != userID {
		t.Errorf("userID incorrect, got: %s, want: %s.", user.UserID, userID)
	}
}
