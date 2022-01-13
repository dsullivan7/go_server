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

func TestIndustries(t *testing.T) {
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

		name := "Name"
		industry := models.Industry{Name: &name}

		db.Create(&industry)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodGet,
			fmt.Sprint(testServer.URL, "/api/industries/", industry.IndustryID),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var industryResponse models.Industry

		errDecode := decoder.Decode(&industryResponse)
		assert.Nil(t, errDecode)

		assert.Equal(t, *industryResponse.Name, *industry.Name)
	})

	t.Run("Test List", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		name1 := "Name1"
		industry1 := models.Industry{Name: &name1}

		name2 := "Name2"
		industry2 := models.Industry{Name: &name2}

		db.Create(&industry1)
		db.Create(&industry2)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodGet,
			fmt.Sprint(testServer.URL, "/api/industries"),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var industriesFound []models.Industry
		errDecode1 := decoder.Decode(&industriesFound)
		assert.Nil(t, errDecode1)

		assert.Equal(t, len(industriesFound), 2)

		var industryResponse models.Industry

		for _, value := range industriesFound {
			if value.IndustryID == industry1.IndustryID {
				industryResponse = value

				break
			}
		}

		assert.Equal(t, industryResponse.IndustryID, industry1.IndustryID)
		assert.Equal(t, *industryResponse.Name, *industry1.Name)

		for _, value := range industriesFound {
			if value.IndustryID == industry2.IndustryID {
				industryResponse = value

				break
			}
		}

		assert.Equal(t, industryResponse.IndustryID, industry2.IndustryID)
		assert.Equal(t, *industryResponse.Name, *industry2.Name)
	})

	t.Run("Test Create", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		jsonStr := []byte(`{"name":"Name"}`)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodPost,
			fmt.Sprint(testServer.URL, "/api/industries"),
			bytes.NewBuffer(jsonStr),
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var industryResponse models.Industry
		errDecode := decoder.Decode(&industryResponse)
		assert.Nil(t, errDecode)

		assert.NotNil(t, industryResponse.IndustryID)
		assert.Equal(t, "Name", *industryResponse.Name)

		var industryFound models.Industry
		errFound := db.Where("industry_id = ?", industryResponse.IndustryID).First(&industryFound).Error

		assert.Nil(t, errFound)

		assert.Equal(t, "Name", *industryFound.Name)
	})

	t.Run("Test Modify", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		name := "Name"
		industry := models.Industry{Name: &name}

		db.Create(&industry)

		jsonStr := []byte(`{"name":"NameDifferent"}`)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodPut,
			fmt.Sprint(testServer.URL, "/api/industries/", industry.IndustryID),
			bytes.NewBuffer(jsonStr),
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var industryResponse models.Industry
		errDecode := decoder.Decode(&industryResponse)
		assert.Nil(t, errDecode)

		assert.Equal(t, industryResponse.IndustryID, industry.IndustryID)
		assert.Equal(t, "NameDifferent", *industryResponse.Name)

		var industryFound models.Industry
		errFound := db.Where("industry_id = ?", industry.IndustryID).First(&industryFound).Error

		assert.Nil(t, errFound)

		assert.Equal(t, industryFound.IndustryID, industry.IndustryID)
		assert.Equal(t, "NameDifferent", *industryFound.Name)
	})

	t.Run("Test Delete", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		name := "Name"
		industry := models.Industry{Name: &name}

		db.Create(&industry)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodDelete,
			fmt.Sprint(testServer.URL, "/api/industries/", industry.IndustryID),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		var industryFound models.Industry
		errFound := db.Where("industry_id = ?", industry.IndustryID).First(&industryFound).Error

		assert.Equal(t, gorm.ErrRecordNotFound, errFound)
	})
}