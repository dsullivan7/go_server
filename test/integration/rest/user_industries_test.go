package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go_server/internal/models"
	"go_server/test/mocks/consts"
	testUtils "go_server/test/utils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserIndustries(t *testing.T) {
	setupUtils := testUtils.NewSetupUtility()

	testServer, db, dbUtility, errIntSetup := setupUtils.SetupIntegration()
	assert.Nil(t, errIntSetup)

	defer testServer.Close()

	context := context.Background()

	t.Run("Test Get", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		user := models.User{}
		db.Create(&user)

		name := "Name"
		industry := models.Industry{Name: &name}
		db.Create(&industry)

		userIndustry := models.UserIndustry{UserID: user.UserID, IndustryID: industry.IndustryID}

		db.Create(&userIndustry)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodGet,
			fmt.Sprint(testServer.URL, "/api/user-industries/", userIndustry.UserIndustryID),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var userIndustryResponse models.UserIndustry

		errDecode := decoder.Decode(&userIndustryResponse)
		assert.Nil(t, errDecode)

		assert.Equal(t, userIndustryResponse.UserID, userIndustry.UserID)
		assert.Equal(t, userIndustryResponse.IndustryID, userIndustry.IndustryID)
	})

	t.Run("Test List", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		user1 := models.User{}
		db.Create(&user1)

		auth0ID := consts.LoggedInAuth0Id

		user2 := models.User{Auth0ID: &auth0ID}
		db.Create(&user2)

		name1 := "Name1"
		industry1 := models.Industry{Name: &name1}
		db.Create(&industry1)

		name2 := "Name2"
		industry2 := models.Industry{Name: &name2}
		db.Create(&industry2)

		userIndustry1 := models.UserIndustry{UserID: user1.UserID, IndustryID: industry1.IndustryID}
		userIndustry2 := models.UserIndustry{UserID: user1.UserID, IndustryID: industry2.IndustryID}
		userIndustry3 := models.UserIndustry{UserID: user2.UserID, IndustryID: industry2.IndustryID}

		db.Create(&userIndustry1)
		db.Create(&userIndustry2)
		db.Create(&userIndustry3)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodGet,
			fmt.Sprint(testServer.URL, "/api/user-industries"),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var userIndustriesFound []models.UserIndustry
		errDecode1 := decoder.Decode(&userIndustriesFound)
		assert.Nil(t, errDecode1)

		assert.Equal(t, 3, len(userIndustriesFound))

		var userIndustryResponse models.UserIndustry

		for _, value := range userIndustriesFound {
			if value.UserIndustryID == userIndustry1.UserIndustryID {
				userIndustryResponse = value

				break
			}
		}

		assert.Equal(t, userIndustryResponse.UserIndustryID, userIndustry1.UserIndustryID)
		assert.Equal(t, userIndustryResponse.UserID, userIndustry1.UserID)
		assert.Equal(t, userIndustryResponse.IndustryID, userIndustry1.IndustryID)

		for _, value := range userIndustriesFound {
			if value.UserIndustryID == userIndustry2.UserIndustryID {
				userIndustryResponse = value

				break
			}
		}

		assert.Equal(t, userIndustryResponse.UserIndustryID, userIndustry2.UserIndustryID)
		assert.Equal(t, userIndustryResponse.UserID, userIndustry2.UserID)
		assert.Equal(t, userIndustryResponse.IndustryID, userIndustry2.IndustryID)

		for _, value := range userIndustriesFound {
			if value.UserIndustryID == userIndustry3.UserIndustryID {
				userIndustryResponse = value

				break
			}
		}

		assert.Equal(t, userIndustryResponse.UserIndustryID, userIndustry3.UserIndustryID)
		assert.Equal(t, userIndustryResponse.UserID, userIndustry3.UserID)
		assert.Equal(t, userIndustryResponse.IndustryID, userIndustry3.IndustryID)

		// test request with query
		req, errRequest = http.NewRequestWithContext(
			context,
			http.MethodGet,
			fmt.Sprint(testServer.URL, "/api/user-industries?user_id=", user2.UserID),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse = http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder = json.NewDecoder(res.Body)

		errDecode2 := decoder.Decode(&userIndustriesFound)
		assert.Nil(t, errDecode2)

		assert.Equal(t, 1, len(userIndustriesFound))

		userIndustryResponse = userIndustriesFound[0]

		assert.Equal(t, userIndustryResponse.UserIndustryID, userIndustry3.UserIndustryID)
		assert.Equal(t, userIndustryResponse.UserID, userIndustry3.UserID)
		assert.Equal(t, userIndustryResponse.IndustryID, userIndustry3.IndustryID)
	})

	t.Run("Test Create", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		user := models.User{}
		db.Create(&user)

		name := "Name"
		industry := models.Industry{Name: &name}
		db.Create(&industry)

		jsonStr := []byte(fmt.Sprintf(`{"user_id":"%s", "industry_id": "%s"}`, user.UserID, industry.IndustryID))

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodPost,
			fmt.Sprint(testServer.URL, "/api/user-industries"),
			bytes.NewBuffer(jsonStr),
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var userIndustryResponse models.UserIndustry
		errDecode := decoder.Decode(&userIndustryResponse)
		assert.Nil(t, errDecode)

		assert.NotNil(t, userIndustryResponse.UserIndustryID)
		assert.Equal(t, user.UserID, userIndustryResponse.UserID)
		assert.Equal(t, industry.IndustryID, userIndustryResponse.IndustryID)

		var userIndustryFound models.UserIndustry
		errFound := db.Where("user_industry_id = ?", userIndustryResponse.UserIndustryID).First(&userIndustryFound).Error

		assert.Nil(t, errFound)

		assert.Equal(t, user.UserID, userIndustryFound.UserID)
		assert.Equal(t, industry.IndustryID, userIndustryFound.IndustryID)
	})

	t.Run("Test Modify", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		user := models.User{}
		db.Create(&user)

		name := "Name"
		industry := models.Industry{Name: &name}
		db.Create(&industry)

		userIndustry := models.UserIndustry{UserID: user.UserID, IndustryID: industry.IndustryID}

		db.Create(&userIndustry)

		jsonStr := []byte(`{}`)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodPut,
			fmt.Sprint(testServer.URL, "/api/user-industries/", userIndustry.UserIndustryID),
			bytes.NewBuffer(jsonStr),
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var userIndustryResponse models.UserIndustry
		errDecode := decoder.Decode(&userIndustryResponse)
		assert.Nil(t, errDecode)

		assert.Equal(t, userIndustryResponse.UserIndustryID, userIndustry.UserIndustryID)
		assert.Equal(t, userIndustryResponse.UserID, userIndustry.UserID)
		assert.Equal(t, userIndustryResponse.IndustryID, userIndustry.IndustryID)

		var userIndustryFound models.UserIndustry
		errFound := db.Where("user_industry_id = ?", userIndustry.UserIndustryID).First(&userIndustryFound).Error

		assert.Nil(t, errFound)

		assert.Equal(t, userIndustryFound.UserIndustryID, userIndustry.UserIndustryID)
		assert.Equal(t, user.UserID, userIndustryFound.UserID)
		assert.Equal(t, industry.IndustryID, userIndustryFound.IndustryID)
	})

	t.Run("Test Delete", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		user := models.User{}
		db.Create(&user)

		name := "Name"
		industry := models.Industry{Name: &name}
		db.Create(&industry)

		userIndustry := models.UserIndustry{UserID: user.UserID, IndustryID: industry.IndustryID}

		db.Create(&userIndustry)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodDelete,
			fmt.Sprint(testServer.URL, "/api/user-industries/", userIndustry.UserIndustryID),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		var userIndustryFound models.UserIndustry
		errFound := db.Where("user_industry_id = ?", userIndustry.UserIndustryID).First(&userIndustryFound).Error

		assert.Equal(t, gorm.ErrRecordNotFound, errFound)
	})
}
