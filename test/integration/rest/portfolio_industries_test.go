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

func TestPortfolioIndustries(t *testing.T) {
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

		portfolio := models.Portfolio{UserID: &user.UserID}
		db.Create(&portfolio)

		name := "Name"
		industry := models.Industry{Name: &name}
		db.Create(&industry)

		portfolioIndustry := models.PortfolioIndustry{PortfolioID: portfolio.PortfolioID, IndustryID: industry.IndustryID}

		db.Create(&portfolioIndustry)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodGet,
			fmt.Sprint(testServer.URL, "/api/portfolio-industries/", portfolioIndustry.PortfolioIndustryID),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var portfolioIndustryResponse models.PortfolioIndustry

		errDecode := decoder.Decode(&portfolioIndustryResponse)
		assert.Nil(t, errDecode)

		assert.Equal(t, portfolioIndustryResponse.PortfolioID, portfolioIndustry.PortfolioID)
		assert.Equal(t, portfolioIndustryResponse.IndustryID, portfolioIndustry.IndustryID)
	})

	t.Run("Test List", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		user1 := models.User{}
		db.Create(&user1)

		portfolio1 := models.Portfolio{UserID: &user1.UserID}
		db.Create(&portfolio1)

		auth0ID := consts.LoggedInAuth0Id

		user2 := models.User{Auth0ID: &auth0ID}
		db.Create(&user2)

		portfolio2 := models.Portfolio{UserID: &user2.UserID}
		db.Create(&portfolio2)

		name1 := "Name1"
		industry1 := models.Industry{Name: &name1}
		db.Create(&industry1)

		name2 := "Name2"
		industry2 := models.Industry{Name: &name2}
		db.Create(&industry2)

		portfolioIndustry1 := models.PortfolioIndustry{PortfolioID: portfolio1.PortfolioID, IndustryID: industry1.IndustryID}
		portfolioIndustry2 := models.PortfolioIndustry{PortfolioID: portfolio1.PortfolioID, IndustryID: industry2.IndustryID}
		portfolioIndustry3 := models.PortfolioIndustry{PortfolioID: portfolio2.PortfolioID, IndustryID: industry2.IndustryID}

		db.Create(&portfolioIndustry1)
		db.Create(&portfolioIndustry2)
		db.Create(&portfolioIndustry3)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodGet,
			fmt.Sprint(testServer.URL, "/api/portfolio-industries"),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var portfolioIndustriesFound []models.PortfolioIndustry
		errDecode1 := decoder.Decode(&portfolioIndustriesFound)
		assert.Nil(t, errDecode1)

		assert.Equal(t, 3, len(portfolioIndustriesFound))

		var portfolioIndustryResponse models.PortfolioIndustry

		for _, value := range portfolioIndustriesFound {
			if value.PortfolioIndustryID == portfolioIndustry1.PortfolioIndustryID {
				portfolioIndustryResponse = value

				break
			}
		}

		assert.Equal(t, portfolioIndustryResponse.PortfolioIndustryID, portfolioIndustry1.PortfolioIndustryID)
		assert.Equal(t, portfolioIndustryResponse.PortfolioID, portfolioIndustry1.PortfolioID)
		assert.Equal(t, portfolioIndustryResponse.IndustryID, portfolioIndustry1.IndustryID)

		for _, value := range portfolioIndustriesFound {
			if value.PortfolioIndustryID == portfolioIndustry2.PortfolioIndustryID {
				portfolioIndustryResponse = value

				break
			}
		}

		assert.Equal(t, portfolioIndustryResponse.PortfolioIndustryID, portfolioIndustry2.PortfolioIndustryID)
		assert.Equal(t, portfolioIndustryResponse.PortfolioID, portfolioIndustry2.PortfolioID)
		assert.Equal(t, portfolioIndustryResponse.IndustryID, portfolioIndustry2.IndustryID)

		for _, value := range portfolioIndustriesFound {
			if value.PortfolioIndustryID == portfolioIndustry3.PortfolioIndustryID {
				portfolioIndustryResponse = value

				break
			}
		}

		assert.Equal(t, portfolioIndustryResponse.PortfolioIndustryID, portfolioIndustry3.PortfolioIndustryID)
		assert.Equal(t, portfolioIndustryResponse.PortfolioID, portfolioIndustry3.PortfolioID)
		assert.Equal(t, portfolioIndustryResponse.IndustryID, portfolioIndustry3.IndustryID)

		// test request with query
		req, errRequest = http.NewRequestWithContext(
			context,
			http.MethodGet,
			fmt.Sprint(testServer.URL, "/api/portfolio-industries?user_id=", user2.UserID),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse = http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder = json.NewDecoder(res.Body)

		errDecode2 := decoder.Decode(&portfolioIndustriesFound)
		assert.Nil(t, errDecode2)

		assert.Equal(t, 1, len(portfolioIndustriesFound))

		portfolioIndustryResponse = portfolioIndustriesFound[0]

		assert.Equal(t, portfolioIndustryResponse.PortfolioIndustryID, portfolioIndustry3.PortfolioIndustryID)
		assert.Equal(t, portfolioIndustryResponse.PortfolioID, portfolioIndustry3.PortfolioID)
		assert.Equal(t, portfolioIndustryResponse.IndustryID, portfolioIndustry3.IndustryID)
	})

	t.Run("Test Create", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		user := models.User{}
		db.Create(&user)

		portfolio := models.Portfolio{UserID: &user.UserID}
		db.Create(&portfolio)

		name := "Name"
		industry := models.Industry{Name: &name}
		db.Create(&industry)

		jsonStr := []byte(fmt.Sprintf(
			`{"portfolio_id":"%s", "industry_id": "%s"}`,
			portfolio.PortfolioID,
			industry.IndustryID,
		))

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodPost,
			fmt.Sprint(testServer.URL, "/api/portfolio-industries"),
			bytes.NewBuffer(jsonStr),
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var portfolioIndustryResponse models.PortfolioIndustry
		errDecode := decoder.Decode(&portfolioIndustryResponse)
		assert.Nil(t, errDecode)

		assert.NotNil(t, portfolioIndustryResponse.PortfolioIndustryID)
		assert.Equal(t, portfolio.PortfolioID, portfolioIndustryResponse.PortfolioID)
		assert.Equal(t, industry.IndustryID, portfolioIndustryResponse.IndustryID)

		var portfolioIndustryFound models.PortfolioIndustry
		errFound := db.Where(
			"portfolio_industry_id = ?",
			portfolioIndustryResponse.PortfolioIndustryID,
		).First(&portfolioIndustryFound).Error

		assert.Nil(t, errFound)

		assert.Equal(t, portfolio.PortfolioID, portfolioIndustryFound.PortfolioID)
		assert.Equal(t, industry.IndustryID, portfolioIndustryFound.IndustryID)
	})

	t.Run("Test Modify", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		user := models.User{}
		db.Create(&user)

		portfolio := models.Portfolio{UserID: &user.UserID}
		db.Create(&portfolio)

		name := "Name"
		industry := models.Industry{Name: &name}
		db.Create(&industry)

		portfolioIndustry := models.PortfolioIndustry{PortfolioID: portfolio.PortfolioID, IndustryID: industry.IndustryID}

		db.Create(&portfolioIndustry)

		jsonStr := []byte(`{}`)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodPut,
			fmt.Sprint(testServer.URL, "/api/portfolio-industries/", portfolioIndustry.PortfolioIndustryID),
			bytes.NewBuffer(jsonStr),
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var portfolioIndustryResponse models.PortfolioIndustry
		errDecode := decoder.Decode(&portfolioIndustryResponse)
		assert.Nil(t, errDecode)

		assert.Equal(t, portfolioIndustryResponse.PortfolioIndustryID, portfolioIndustry.PortfolioIndustryID)
		assert.Equal(t, portfolioIndustryResponse.PortfolioID, portfolioIndustry.PortfolioID)
		assert.Equal(t, portfolioIndustryResponse.IndustryID, portfolioIndustry.IndustryID)

		var portfolioIndustryFound models.PortfolioIndustry
		errFound := db.Where(
			"portfolio_industry_id = ?",
			portfolioIndustry.PortfolioIndustryID,
		).First(&portfolioIndustryFound).Error

		assert.Nil(t, errFound)

		assert.Equal(t, portfolioIndustryFound.PortfolioIndustryID, portfolioIndustry.PortfolioIndustryID)
		assert.Equal(t, portfolio.PortfolioID, portfolioIndustryFound.PortfolioID)
		assert.Equal(t, industry.IndustryID, portfolioIndustryFound.IndustryID)
	})

	t.Run("Test Delete", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		user := models.User{}
		db.Create(&user)

		portfolio := models.Portfolio{UserID: &user.UserID}
		db.Create(&portfolio)

		name := "Name"
		industry := models.Industry{Name: &name}
		db.Create(&industry)

		portfolioIndustry := models.PortfolioIndustry{PortfolioID: portfolio.PortfolioID, IndustryID: industry.IndustryID}

		db.Create(&portfolioIndustry)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodDelete,
			fmt.Sprint(testServer.URL, "/api/portfolio-industries/", portfolioIndustry.PortfolioIndustryID),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		var portfolioIndustryFound models.PortfolioIndustry
		errFound := db.Where(
			"portfolio_industry_id = ?",
			portfolioIndustry.PortfolioIndustryID,
		).First(&portfolioIndustryFound).Error

		assert.Equal(t, gorm.ErrRecordNotFound, errFound)
	})
}
