package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go_server/internal/models"
	testUtils "go_server/test/utils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UsersResponse struct {
	Data struct {
		Users []models.User `json:"users"`
	} `json:"data"`
}

func TestUsers(t *testing.T) {
	setupUtils := testUtils.NewSetupUtility()

	testServer, db, dbUtility, errIntSetup := setupUtils.SetupIntegration()
	assert.Nil(t, errIntSetup)

	defer testServer.Close()

	context := context.Background()

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
