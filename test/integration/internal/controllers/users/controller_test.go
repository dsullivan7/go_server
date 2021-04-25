package users_test

import (
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"errors"

	"gorm.io/gorm"

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
		t.Fatalf("Expected %d, Received %d", http.StatusOK, res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)

	var userResponse models.User
	errDecode := decoder.Decode(&userResponse)

	if errDecode != nil {
		t.Fatalf("Decoding error: %v", errDecode)
	}

	if user.UserID != userResponse.UserID {
		t.Fatalf("Expected: %s, Received: %s", user.UserID, userResponse.UserID)
	}

	if user.FirstName != userResponse.FirstName {
		t.Fatalf("Expected: %s, Received: %s", user.FirstName, userResponse.FirstName)
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
		t.Fatalf("Expected %d, Received %d", http.StatusOK, res.StatusCode)
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

	var userResponse models.User

	for _, value := range usersFound {
    if value.UserID == user1.UserID {
			userResponse = value
			break
    }
	}

	if user1.FirstName != userResponse.FirstName {
		t.Fatalf("Expected: %s, Received: %s", user1.FirstName, userResponse.FirstName)
	}

	for _, value := range usersFound {
    if value.UserID == user2.UserID {
			userResponse = value
			break
    }
	}

	if user2.FirstName != userResponse.FirstName {
		t.Fatalf("Expected: %s, Received: %s", user2.FirstName, userResponse.FirstName)
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

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Expected %d, Received %d", http.StatusCreated, res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)

	var userResponse models.User
	errDecode := decoder.Decode(&userResponse)

	if errDecode != nil {
		t.Fatalf("Decoding error: %v", errDecode)
	}

	if userResponse.FirstName != "FirstName" {
		t.Fatalf("Expected: %s, Received: %s", "FirstName", userResponse.FirstName)
	}

	var userFound models.User
	errFound := db.DB.Where("user_id = ?", userResponse.UserID).First(&userFound).Error

	if errFound != nil {
		t.Fatalf("Error: %v", errFound)
	}

	if "FirstName" != userFound.FirstName {
		t.Fatalf("Expected: %s, Received: %s", "FirstName", userFound.FirstName)
	}
}

func TestModify(t *testing.T) {
	db.Connect()
	db.DB.Exec("truncate table users")

	testServer := httptest.NewServer(routes.Init())
	defer testServer.Close()

  user := models.User{ FirstName: "FirstName" }

  db.DB.Create(&user)

	userDifferent := models.User{ FirstName: "FirstNameDifferent" }
	jsonReq, errJSON := json.Marshal(userDifferent)

	if errJSON != nil {
		t.Fatalf("JSON: %v", errJSON)
	}

	req, errRequest := http.NewRequest(http.MethodPut, fmt.Sprint(testServer.URL, "/api/users/", user.UserID), bytes.NewBuffer(jsonReq))
	client := &http.Client{}
	res, errRequest := client.Do(req)

	if errRequest != nil {
		t.Fatalf("Delete: %v", errRequest)
	}

	if res.StatusCode != http.StatusOK {
    t.Fatalf("Expected %d, Received %d", http.StatusOK, res.StatusCode)
  }

	decoder := json.NewDecoder(res.Body)

	var userResponse models.User
	errDecode := decoder.Decode(&userResponse)

	if errDecode != nil {
		t.Fatalf("Decoding error: %v", errDecode)
	}

	if user.UserID != userResponse.UserID {
		t.Fatalf("Expected: %s, Received: %s", user.UserID, userResponse.UserID)
	}

	if userDifferent.FirstName != userResponse.FirstName {
		t.Fatalf("Expected: %s, Received: %s", user.FirstName, userResponse.FirstName)
	}

	var userFound models.User
  errFound := db.DB.Where("user_id = ?", user.UserID).First(&userFound).Error

	if errFound != nil {
		t.Fatalf("Error: %v", errFound)
	}

	if user.UserID != userFound.UserID {
		t.Fatalf("Expected: %s, Received: %s", user.UserID, userFound.UserID)
	}

	if userDifferent.FirstName != userFound.FirstName {
		t.Fatalf("Expected: %s, Received: %s", user.FirstName, userFound.FirstName)
	}
}

func TestDelete(t *testing.T) {
	db.Connect()
	db.DB.Exec("truncate table users")

	testServer := httptest.NewServer(routes.Init())
	defer testServer.Close()

  user := models.User{ FirstName: "FirstName" }

  db.DB.Create(&user)

	req, errRequest := http.NewRequest(http.MethodDelete, fmt.Sprint(testServer.URL, "/api/users/", user.UserID), nil)
	client := &http.Client{}
	res, errRequest := client.Do(req)

	if errRequest != nil {
		t.Fatalf("Delete: %v", errRequest)
	}

	if res.StatusCode != http.StatusNoContent {
    t.Fatalf("Expected %d, Received %d", http.StatusNoContent, res.StatusCode)
  }

	var userFound models.User
  errFound := db.DB.Where("user_id = ?", user.UserID).First(&userFound).Error

	if !errors.Is(errFound, gorm.ErrRecordNotFound) {
		t.Fatalf("Expected user not to be found")
	}
}
