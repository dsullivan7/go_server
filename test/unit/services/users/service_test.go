package users

import (
  "testing"
  "github.com/google/uuid"

  UsersService "go_server/internal/services/users"
)

func TestGet(t *testing.T) {
    ID := uuid.New()
    user := UsersService.Get(ID)
    if user.UserID != ID {
       t.Errorf("userID incorrect, got: %s, want: %s.", user.UserID, ID)
    }
}
