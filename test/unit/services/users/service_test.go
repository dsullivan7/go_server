package users_test

import (
	"testing"

	UsersService "go_server/internal/services/users"

	"github.com/google/uuid"
)

func TestGet(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	user := UsersService.Get(userID)

	if user.UserID != userID {
		t.Errorf("userID incorrect, got: %s, want: %s.", user.UserID, userID)
	}
}
