package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go_server/internal/captcha/twocaptcha"
	"go_server/internal/config"
	goServerRodCrawler "go_server/internal/crawler/rod"
	"go_server/internal/db"
	goServerZapLogger "go_server/internal/logger/zap"
	"go_server/internal/models"
	"go_server/internal/server"
	goServerGormStore "go_server/internal/store/gorm"
	"go_server/test/mocks/auth"
	"go_server/test/mocks/consts"
	testUtils "go_server/test/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-rod/rod"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestUsers(t *testing.T) {
	config, configError := config.NewConfig()
	assert.Nil(t, configError)

	zapLogger, errZap := zap.NewProduction()
	assert.Nil(t, errZap)

	logger := goServerZapLogger.NewLogger(zapLogger)

	connection, errConnection := db.NewSQLConnection(
		config.DBHost,
		config.DBName,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBSSL,
	)
	assert.Nil(t, errConnection)

	db, errDatabase := db.NewGormDB(connection)
	assert.Nil(t, errDatabase)

	dbUtility := testUtils.NewSQLDatabaseUtility(connection)

	store := goServerGormStore.NewStore(db)

	router := chi.NewRouter()

	authMock := auth.NewMockAuth()

	browser := rod.New()

	captchaKey := "key"

	captcha := twocaptcha.NewTwoCaptcha(captchaKey, logger)

	crawler := goServerRodCrawler.NewCrawler(browser, captcha)

	handler := server.NewChiServer(config, router, store, crawler, authMock, logger)

	testServer := httptest.NewServer(handler.Init())

	context := context.Background()

	defer testServer.Close()

	t.Run("Test Get", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		firstName := "firstName"
		lastName := "lastName"
		auth0ID := consts.LoggedInAuth0Id

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

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var userResponse models.User
		errDecoder := decoder.Decode(&userResponse)
		assert.Nil(t, errDecoder)

		assert.Equal(t, userResponse.UserID, user.UserID)
		assert.Equal(t, *userResponse.FirstName, *user.FirstName)
		assert.Equal(t, *userResponse.LastName, *user.LastName)
		assert.Equal(t, *userResponse.Auth0ID, *user.Auth0ID)
	})

	t.Run("Test Get Me", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		firstName := "firstName"
		lastName := "lastName"
		auth0ID := consts.LoggedInAuth0Id

		user := models.User{
			FirstName: &firstName,
			LastName:  &lastName,
			Auth0ID:   &auth0ID,
		}

		db.Create(&user)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodGet,
			fmt.Sprint(testServer.URL, "/api/users/me"),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var userResponse models.User
		errDecoder := decoder.Decode(&userResponse)
		assert.Nil(t, errDecoder)

		assert.Equal(t, userResponse.UserID, user.UserID)
		assert.Equal(t, *userResponse.FirstName, *user.FirstName)
		assert.Equal(t, *userResponse.LastName, *user.LastName)
		assert.Equal(t, *userResponse.Auth0ID, *user.Auth0ID)
	})

	t.Run("Test List", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

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

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var usersFound []models.User
		errDecoder := decoder.Decode(&usersFound)
		assert.Nil(t, errDecoder)

		assert.Equal(t, 2, len(usersFound))

		var userResponse models.User

		for _, value := range usersFound {
			if value.UserID == user1.UserID {
				userResponse = value

				break
			}
		}

		assert.Equal(t, userResponse.UserID, user1.UserID)
		assert.Equal(t, *userResponse.FirstName, *user1.FirstName)
		assert.Equal(t, *userResponse.LastName, *user1.LastName)
		assert.Equal(t, *userResponse.Auth0ID, *user1.Auth0ID)

		for _, value := range usersFound {
			if value.UserID == user2.UserID {
				userResponse = value

				break
			}
		}

		assert.Equal(t, userResponse.UserID, user2.UserID)
		assert.Equal(t, *userResponse.FirstName, *user2.FirstName)
		assert.Equal(t, *userResponse.LastName, *user2.LastName)
		assert.Equal(t, *userResponse.Auth0ID, *user2.Auth0ID)
	})

	t.Run("Test Create", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

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

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var userResponse models.User
		errDecoder := decoder.Decode(&userResponse)
		assert.Nil(t, errDecoder)

		assert.NotNil(t, userResponse.UserID)
		assert.Equal(t, "FirstName", *userResponse.FirstName)
		assert.Equal(t, "LastName", *userResponse.LastName)
		// assert.Equal(t, *userResponse.Auth0ID, "auth0|loggedInUser")

		var userFound models.User
		errFound := db.Where("user_id = ?", userResponse.UserID).First(&userFound).Error

		assert.Nil(t, errFound)

		assert.Equal(t, "FirstName", *userFound.FirstName)
		assert.Equal(t, "LastName", *userFound.LastName)
		// assert.Equal(t, *userFound.Auth0ID, "auth0|loggedInUser")
	})

	t.Run("Test Modify", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

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

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var userResponse models.User
		errDecoder := decoder.Decode(&userResponse)
		assert.Nil(t, errDecoder)

		assert.Equal(t, "FirstNameDifferent", *userResponse.FirstName)
		assert.Equal(t, "LastNameDifferent", *userResponse.LastName)
		assert.Equal(t, "Auth0IDDifferent", *userResponse.Auth0ID)

		var userFound models.User
		errFound := db.Where("user_id = ?", user.UserID).First(&userFound).Error

		assert.Nil(t, errFound)

		assert.Equal(t, "FirstNameDifferent", *userFound.FirstName)
		assert.Equal(t, "LastNameDifferent", *userFound.LastName)
		assert.Equal(t, "Auth0IDDifferent", *userFound.Auth0ID)
	})

	t.Run("Test Delete", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

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

		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		var userFound models.User
		errFound := db.Where("user_id = ?", user.UserID).First(&userFound).Error

		assert.Equal(t, errFound, gorm.ErrRecordNotFound)
	})
}
