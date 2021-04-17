package users_test

import (
	"fmt"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go_server/internal/routes"
	"go_server/internal/models"
	"go_server/internal/db"
)

func TestGet(t *testing.T) {
	t.Parallel()

	testServer := httptest.NewServer(routes.Init())
	defer testServer.Close()

  user := models.User{ FirstName: "FirstName" }
  db.Connect()
  db.DB.Create(&user)

	res, errRequest := http.Get(fmt.Sprint(testServer.URL, "/api/users/", user.UserID))

	if errRequest != nil {
		t.Fatalf("Get: %v", errRequest)
	}

	if res.StatusCode != http.StatusOK {
    t.Fatalf("didn’t respond 200 OK: %s", res.Status)
  }

	decoder := json.NewDecoder(res.Body)

	var userFound models.User
	errDecode := decoder.Decode(&userFound)

	if errDecode != nil {
		t.Fatalf("Decoding error: %v", errDecode)
	}

	if user.UserID != userFound.UserID {
		t.Fatalf("Expected: %s, Received: %s", user.UserID, userFound.UserID)
	}

	if user.FirstName != userFound.FirstName {
		t.Fatalf("Expected: %s, Received: %s", user.FirstName, userFound.FirstName)
	}
}
