package integration_test

import (
	// jwt "github.com/dgrijalva/jwt-go".
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go_server/internal/config"
	"go_server/internal/controllers"
	"go_server/internal/db"
	"go_server/internal/logger"
	"go_server/internal/models"
	"go_server/internal/server"
	"go_server/internal/store"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// type userString string
// const userKey = userString("user")

// func init() {
// middlewares.Auth = func(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		jwtToken := &jwt.Token{ Claims: jwt.MapClaims{ "sub": "auth0|loggedInUser" } }
//     newContext := context.WithValue(r.Context(), userKey, jwtToken)
//     h.ServeHTTP(w, r.WithContext(newContext))
//   })
// }
// }

func TestUsers(t *testing.T) {
	config, configError := config.NewConfig()
	assert.Nil(t, configError)

	zapLogger, errZap := zap.NewProduction()
	assert.Nil(t, errZap)

	logger := logger.NewZapLogger(zapLogger)

	db, errDatabase := db.NewDatabase(
		config.DBHost,
		config.DBName,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBSSL,
	)
	assert.Nil(t, errDatabase)

	store := store.NewGormStore(db)

	controllers := controllers.NewControllers(store, config, logger)
	router := chi.NewRouter()
	server := server.NewServer(router, controllers, config, logger)

	testServer := httptest.NewServer(server.Routes())
	context := context.Background()

	defer testServer.Close()

	t.Run("Test Get", func(t *testing.T) {
		store.TruncateAll()

		firstName := "firstName"
		lastName := "lastName"
		auth0ID := "auth0ID"

		user := models.User{
			FirstName: &firstName,
			LastName:  &lastName,
			Auth0ID:   &auth0ID,
		}

		db.Create(&user)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodGet,
			fmt.Sprint(testServer.URL, "/api/users/", user.UserID),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusOK)

		decoder := json.NewDecoder(res.Body)

		var userResponse models.User
		errDecoder := decoder.Decode(&userResponse)
		assert.Nil(t, errDecoder)

		assert.Equal(t, user.UserID, userResponse.UserID)
		assert.Equal(t, *user.FirstName, *userResponse.FirstName)
		assert.Equal(t, *user.LastName, *userResponse.LastName)
		assert.Equal(t, *user.Auth0ID, *userResponse.Auth0ID)
	})

	t.Run("Test List", func(t *testing.T) {
		store.TruncateAll()

		firstName1 := "firstName1"
		lastName1 := "lastName1"
		auth0Id1 := "auth0Id1"

		firstName2 := "firstName2"
		lastName2 := "lastName2"
		auth0Id2 := "auth0Id2"

		user1 := models.User{
			FirstName: &firstName1,
			LastName:  &lastName1,
			Auth0ID:   &auth0Id1,
		}

		user2 := models.User{
			FirstName: &firstName2,
			LastName:  &lastName2,
			Auth0ID:   &auth0Id2,
		}

		db.Create(&user1)
		db.Create(&user2)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodGet,
			fmt.Sprint(testServer.URL, "/api/users"),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusOK)

		decoder := json.NewDecoder(res.Body)

		var usersFound []models.User
		errDecoder := decoder.Decode(&usersFound)
		assert.Nil(t, errDecoder)

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
	})

	t.Run("Test Create", func(t *testing.T) {
		store.TruncateAll()

		jsonStr := []byte(`{
			"first_name":"FirstName",
			"last_name":"LastName"
		}`)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodPost,
			fmt.Sprint(testServer.URL, "/api/users"),
			bytes.NewBuffer(jsonStr),
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusCreated)

		decoder := json.NewDecoder(res.Body)

		var userResponse models.User
		errDecoder := decoder.Decode(&userResponse)
		assert.Nil(t, errDecoder)

		assert.NotNil(t, userResponse.UserID)
		assert.Equal(t, *userResponse.FirstName, "FirstName")
		assert.Equal(t, *userResponse.LastName, "LastName")
		// assert.Equal(t, *userResponse.Auth0ID, "auth0|loggedInUser")

		var userFound models.User
		errFound := db.Where("user_id = ?", userResponse.UserID).First(&userFound).Error

		assert.Nil(t, errFound)

		assert.Equal(t, *userFound.FirstName, "FirstName")
		assert.Equal(t, *userFound.LastName, "LastName")
		// assert.Equal(t, *userFound.Auth0ID, "auth0|loggedInUser")
	})

	t.Run("Test Modify", func(t *testing.T) {
		store.TruncateAll()

		firstName := "FirstName"
		lastName := "LastName"
		auth0ID := "Auth0ID"
		user := models.User{
			FirstName: &firstName,
			LastName:  &lastName,
			Auth0ID:   &auth0ID,
		}

		db.Create(&user)

		jsonStr := []byte(`{
			"first_name":"FirstNameDifferent",
			"last_name": "LastNameDifferent",
			"auth0_id": "Auth0IDDifferent"
		}`)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodPut,
			fmt.Sprint(testServer.URL, "/api/users/", user.UserID),
			bytes.NewBuffer(jsonStr),
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusOK)

		decoder := json.NewDecoder(res.Body)

		var userResponse models.User
		errDecoder := decoder.Decode(&userResponse)
		assert.Nil(t, errDecoder)

		assert.Equal(t, *userResponse.FirstName, "FirstNameDifferent")
		assert.Equal(t, *userResponse.LastName, "LastNameDifferent")
		assert.Equal(t, *userResponse.Auth0ID, "Auth0IDDifferent")

		var userFound models.User
		errFound := db.Where("user_id = ?", user.UserID).First(&userFound).Error

		assert.Nil(t, errFound)

		assert.Equal(t, *userFound.FirstName, "FirstNameDifferent")
		assert.Equal(t, *userFound.LastName, "LastNameDifferent")
		assert.Equal(t, *userFound.Auth0ID, "Auth0IDDifferent")
	})

	t.Run("Test Delete", func(t *testing.T) {
		store.TruncateAll()

		firstName := "firstName"
		user := models.User{FirstName: &firstName}

		db.Create(&user)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodDelete,
			fmt.Sprint(testServer.URL, "/api/users/", user.UserID),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusNoContent)

		var userFound models.User
		errFound := db.Where("user_id = ?", user.UserID).First(&userFound).Error

		assert.Equal(t, errFound, gorm.ErrRecordNotFound)
	})
}
