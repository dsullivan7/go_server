package controllers

import (
	"fmt"
	"context"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"

	"gorm.io/gorm"

	"go_server/internal/routes"
	"go_server/internal/models"
	"go_server/internal/db"
	"go_server/internal/middlewares"

	"github.com/stretchr/testify/assert"
)

func init() {
	middlewares.Auth = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwtToken := &jwt.Token{ Claims: jwt.MapClaims{ "sub": "auth0|loggedInUser" } }
	    newContext := context.WithValue(r.Context(), "user", jwtToken)
      h.ServeHTTP(w, r.WithContext(newContext))
	  })
	}
}

func TestGetUser(t *testing.T) {
	db.Connect()
	db.DB.Exec("truncate table users cascade")

	testServer := httptest.NewServer(routes.Init())
	defer testServer.Close()

	firstName := "firstName"
	lastName := "lastName"
	auth0ID := "auth0ID"

  user := models.User{
		FirstName: &firstName,
		LastName: &lastName,
		Auth0ID: &auth0ID,
	 }

  db.DB.Create(&user)

	res, errRequest := http.Get(fmt.Sprint(testServer.URL, "/api/users/", user.UserID))

	assert.Nil(t, errRequest)

	assert.Equal(t, res.StatusCode, http.StatusOK)

	decoder := json.NewDecoder(res.Body)

	var userResponse models.User
	decoder.Decode(&userResponse)

	assert.Equal(t, user.UserID, userResponse.UserID)
	assert.Equal(t, *user.FirstName, *userResponse.FirstName)
	assert.Equal(t, *user.LastName, *userResponse.LastName)
	assert.Equal(t, *user.Auth0ID, *userResponse.Auth0ID)
}

func TestListUsers(t *testing.T) {
	db.Connect()
	db.DB.Exec("truncate table users cascade")

	testServer := httptest.NewServer(routes.Init())
	defer testServer.Close()

	firstName1 := "firstName1"
	lastName1 := "lastName1"
	auth0Id1 := "auth0Id1"

	firstName2 := "firstName2"
	lastName2 := "lastName2"
	auth0Id2 := "auth0Id2"

  user1 := models.User{
		FirstName: &firstName1,
		LastName: &lastName1,
		Auth0ID: &auth0Id1,
	}

  user2 := models.User{
		FirstName: &firstName2,
		LastName: &lastName2,
		Auth0ID: &auth0Id2,
	}

  db.DB.Create(&user1)
	db.DB.Create(&user2)

	res, errRequest := http.Get(fmt.Sprint(testServer.URL, "/api/users"))

	assert.Nil(t, errRequest)

	assert.Equal(t, res.StatusCode, http.StatusOK)

	decoder := json.NewDecoder(res.Body)

	var usersFound []models.User
	decoder.Decode(&usersFound)

	assert.Equal(t, len(usersFound), 2)

	var userResponse models.User

	for _, value := range usersFound {
    if value.UserID == user1.UserID {
			userResponse = value
			break
    }
	}

	assert.Equal(t, user1.UserID, userResponse.UserID)
	assert.Equal(t, *user1.FirstName, *userResponse.FirstName)
	assert.Equal(t, *user1.LastName, *userResponse.LastName)
	assert.Equal(t, *user1.Auth0ID, *userResponse.Auth0ID)

	for _, value := range usersFound {
    if value.UserID == user2.UserID {
			userResponse = value
			break
    }
	}

	assert.Equal(t, user2.UserID, userResponse.UserID)
	assert.Equal(t, *user2.FirstName, *userResponse.FirstName)
	assert.Equal(t, *user2.LastName, *userResponse.LastName)
	assert.Equal(t, *user2.Auth0ID, *userResponse.Auth0ID)
}

func TestCreateUser(t *testing.T) {
	db.Connect()
	db.DB.Exec("truncate table users cascade")

	testServer := httptest.NewServer(routes.Init())
	defer testServer.Close()

	var jsonStr = []byte(`{
		"first_name":"FirstName",
		"last_name":"LastName"
	}`)

	res, errRequest := http.Post(fmt.Sprint(testServer.URL, "/api/users"), "application/json", bytes.NewBuffer(jsonStr))

	assert.Nil(t, errRequest)

	assert.Equal(t, res.StatusCode, http.StatusCreated)

	decoder := json.NewDecoder(res.Body)

	var userResponse models.User
	decoder.Decode(&userResponse)

	assert.NotNil(t, userResponse.UserID)
	assert.Equal(t, *userResponse.FirstName, "FirstName")
	assert.Equal(t, *userResponse.LastName, "LastName")
	// assert.Equal(t, *userResponse.Auth0ID, "auth0|loggedInUser")

	var userFound models.User
	errFound := db.DB.Where("user_id = ?", userResponse.UserID).First(&userFound).Error

	assert.Nil(t, errFound)

	assert.Equal(t, *userFound.FirstName, "FirstName")
	assert.Equal(t, *userFound.LastName, "LastName")
	// assert.Equal(t, *userFound.Auth0ID, "auth0|loggedInUser")
}

func TestModifyUser(t *testing.T) {
	db.Connect()
	db.DB.Exec("truncate table users cascade")

	testServer := httptest.NewServer(routes.Init())
	defer testServer.Close()

	firstName:= "FirstName"
	lastName:= "LastName"
	auth0ID:= "Auth0ID"
  user := models.User{
		FirstName: &firstName,
		LastName: &lastName,
		Auth0ID: &auth0ID,
	}

  db.DB.Create(&user)

	var jsonStr = []byte(`{
		"first_name":"FirstNameDifferent",
		"last_name": "LastNameDifferent",
		"auth0_id": "Auth0IDDifferent"
	}`)

	req, errRequest := http.NewRequest(http.MethodPut, fmt.Sprint(testServer.URL, "/api/users/", user.UserID), bytes.NewBuffer(jsonStr))
	client := &http.Client{}
	res, errRequest := client.Do(req)

	assert.Nil(t, errRequest)

	assert.Equal(t, res.StatusCode, http.StatusOK)

	decoder := json.NewDecoder(res.Body)

	var userResponse models.User
	decoder.Decode(&userResponse)

	assert.Equal(t, *userResponse.FirstName, "FirstNameDifferent")
	assert.Equal(t, *userResponse.LastName, "LastNameDifferent")
	assert.Equal(t, *userResponse.Auth0ID, "Auth0IDDifferent")

	var userFound models.User
  errFound := db.DB.Where("user_id = ?", user.UserID).First(&userFound).Error

	assert.Nil(t, errFound)

	assert.Equal(t, *userFound.FirstName, "FirstNameDifferent")
	assert.Equal(t, *userFound.LastName, "LastNameDifferent")
	assert.Equal(t, *userFound.Auth0ID, "Auth0IDDifferent")
}

func TestDeleteUser(t *testing.T) {
	db.Connect()
	db.DB.Exec("truncate table users cascade")

	testServer := httptest.NewServer(routes.Init())
	defer testServer.Close()

	firstName := "firstName"
  user := models.User{ FirstName: &firstName }

  db.DB.Create(&user)

	req, errRequest := http.NewRequest(http.MethodDelete, fmt.Sprint(testServer.URL, "/api/users/", user.UserID), nil)
	client := &http.Client{}
	res, errRequest := client.Do(req)

	assert.Nil(t, errRequest)

	assert.Equal(t, res.StatusCode, http.StatusNoContent)

	var userFound models.User
  errFound := db.DB.Where("user_id = ?", user.UserID).First(&userFound).Error

	assert.Equal(t, errFound, gorm.ErrRecordNotFound)
}
