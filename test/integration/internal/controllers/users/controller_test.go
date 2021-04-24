package users_test

import (
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go_server/internal/routes"
	"go_server/internal/models"
	"go_server/internal/db"
)

func TestGet(t *testing.T) {
	db.Connect()
	db.DB.Exec("truncate table users")

	testServer := httptest.NewServer(routes.Init())
	defer testServer.Close()

  user := models.User{ FirstName: "FirstName" }

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

func TestList(t *testing.T) {
	db.Connect()
	db.DB.Exec("truncate table users")

	testServer := httptest.NewServer(routes.Init())
	defer testServer.Close()

  user1 := models.User{ FirstName: "FirstName1" }
  user2 := models.User{ FirstName: "FirstName2" }
  db.DB.Create(&user1)
	db.DB.Create(&user2)

	res, errRequest := http.Get(fmt.Sprint(testServer.URL, "/api/users"))

	if errRequest != nil {
		t.Fatalf("Get: %v", errRequest)
	}

	if res.StatusCode != http.StatusOK {
    t.Fatalf("didn’t respond 201 OK: %s", res.Status)
  }

	decoder := json.NewDecoder(res.Body)

	var usersFound []models.User
	errDecode := decoder.Decode(&usersFound)

	if errDecode != nil {
		t.Fatalf("Decoding error: %v", errDecode)
	}

	if len(usersFound) != 2 {
		t.Fatalf("Expected: %d, Received: %d", 2, len(usersFound))
	}

	var userFound models.User

	for _, value := range usersFound {
    if value.UserID == user1.UserID {
			userFound = value
			break
    }
	}

	if user1.FirstName != userFound.FirstName {
		t.Fatalf("Expected: %s, Received: %s", user1.FirstName, userFound.FirstName)
	}

	for _, value := range usersFound {
    if value.UserID == user2.UserID {
			userFound = value
			break
    }
	}

	if user2.FirstName != userFound.FirstName {
		t.Fatalf("Expected: %s, Received: %s", user2.FirstName, userFound.FirstName)
	}
}

func TestCreate(t *testing.T) {
	db.Connect()
	db.DB.Exec("truncate table users")

	testServer := httptest.NewServer(routes.Init())
	defer testServer.Close()

	var jsonStr = []byte(`{"first_name":"FirstName"}`)

	res, errRequest := http.Post(fmt.Sprint(testServer.URL, "/api/users"), "application/json", bytes.NewBuffer(jsonStr))

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


	if userFound.FirstName != "FirstName" {
		t.Fatalf("Expected: %s, Received: %s", "FirstName", userFound.FirstName)
	}
}
