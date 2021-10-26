package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go_server/internal/config"
	"go_server/internal/db"
	goServerZapLogger "go_server/internal/logger/zap"
	"go_server/internal/models"
	"go_server/internal/server"
	goServerGormStore "go_server/internal/store/gorm"
	"go_server/test/mocks/auth"
	testUtils "go_server/test/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type UsersResponse struct {
	Data struct {
		Users []models.User `json:"users"`
	} `json:"data"`
}

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

	authMock := auth.NewAuth()

	handler := server.NewChiServer(config, router, store, authMock, logger)

	testServer := httptest.NewServer(handler.Init())

	context := context.Background()

	defer testServer.Close()

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

		jsonData := map[string]string{
			"query": `
            {
                users {
                    user_id,
                    first_name,
                    last_name,
                    auth0_id,
                }
            }
        `,
		}

		jsonValue, errMashal := json.Marshal(jsonData)
		assert.Nil(t, errMashal)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodPost,
			fmt.Sprint(testServer.URL, "/query"),
			bytes.NewBuffer(jsonValue),
		)

		req.Header.Add("Content-Type", "application/json")

		res, errResponse := http.DefaultClient.Do(req)
		assert.Nil(t, errRequest)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var userResponse UsersResponse
		errDecoder := decoder.Decode(&userResponse)
		assert.Nil(t, errDecoder)

		assert.Equal(t, 2, len(userResponse.Data.Users))

		var userMatch models.User

		for _, value := range userResponse.Data.Users {
			if value.UserID == user1.UserID {
				userMatch = value

				break
			}
		}

		assert.Equal(t, userMatch.UserID, user1.UserID)
		assert.Equal(t, *userMatch.FirstName, *user1.FirstName)
		assert.Equal(t, *userMatch.LastName, *user1.LastName)
		assert.Equal(t, *userMatch.Auth0ID, *user1.Auth0ID)

		for _, value := range userResponse.Data.Users {
			if value.UserID == user2.UserID {
				userMatch = value

				break
			}
		}

		assert.Equal(t, userMatch.UserID, user2.UserID)
		assert.Equal(t, *userMatch.FirstName, *user2.FirstName)
		assert.Equal(t, *userMatch.LastName, *user2.LastName)
		assert.Equal(t, *userMatch.Auth0ID, *user2.Auth0ID)
	})
}
