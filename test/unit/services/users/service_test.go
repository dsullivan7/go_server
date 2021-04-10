package users

import (
  "testing"
  "github.com/satori/go.uuid"

  UsersService "go_server/internal/services/users"
)

func TestGet(t *testing.T) {
    ID, err := uuid.NewV4()
    user := UsersService.Get(ID)
    if user.UserID != ID {
       t.Errorf("userID incorrect, got: %s, want: %s.", user.UserID, ID)
    }
}
