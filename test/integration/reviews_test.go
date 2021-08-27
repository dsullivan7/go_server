package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go_server/internal/config"
	"go_server/internal/db"
	"go_server/internal/models"
	"go_server/internal/server"
	goServerGormStore "go_server/internal/store/gorm"
	goServerZapLogger "go_server/internal/logger/zap"
	testUtils "go_server/test/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestReviews(t *testing.T) {
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

	handler := server.NewChiServer(config, router, store, logger)

	testServer := httptest.NewServer(handler.Init())

	context := context.Background()

	defer testServer.Close()

	t.Run("Test Get", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		user1 := models.User{}
		user2 := models.User{}

		db.Create(&user1)
		db.Create(&user2)

		text := "Text"
		review := models.Review{FromUserID: &user1.UserID, ToUserID: &user2.UserID, Text: &text}

		db.Create(&review)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodGet,
			fmt.Sprint(testServer.URL, "/api/reviews/", review.ReviewID),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var reviewResponse models.Review

		errDecode := decoder.Decode(&reviewResponse)
		assert.Nil(t, errDecode)

		assert.Equal(t, reviewResponse.ReviewID, review.ReviewID)
		assert.Equal(t, *reviewResponse.Text, *review.Text)
		assert.Equal(t, *reviewResponse.FromUserID, *review.FromUserID)
		assert.Equal(t, *reviewResponse.ToUserID, *review.ToUserID)
	})

	t.Run("Test List", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		user1 := models.User{}
		user2 := models.User{}
		user3 := models.User{}
		user4 := models.User{}

		db.Create(&user1)
		db.Create(&user2)
		db.Create(&user3)
		db.Create(&user4)

		text1 := "Text1"
		review1 := models.Review{FromUserID: &user1.UserID, ToUserID: &user2.UserID, Text: &text1}

		text2 := "Text2"
		review2 := models.Review{FromUserID: &user2.UserID, ToUserID: &user3.UserID, Text: &text2}

		db.Create(&review1)
		db.Create(&review2)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodGet,
			fmt.Sprint(testServer.URL, "/api/reviews"),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var reviewsFound []models.Review
		errDecode1 := decoder.Decode(&reviewsFound)
		assert.Nil(t, errDecode1)

		assert.Equal(t, len(reviewsFound), 2)

		var reviewResponse models.Review

		for _, value := range reviewsFound {
			if value.ReviewID == review1.ReviewID {
				reviewResponse = value

				break
			}
		}

		assert.Equal(t, reviewResponse.ReviewID, review1.ReviewID)
		assert.Equal(t, *reviewResponse.FromUserID, *review1.FromUserID)
		assert.Equal(t, *reviewResponse.ToUserID, *review1.ToUserID)
		assert.Equal(t, *reviewResponse.Text, *review1.Text)

		for _, value := range reviewsFound {
			if value.ReviewID == review2.ReviewID {
				reviewResponse = value

				break
			}
		}

		assert.Equal(t, reviewResponse.ReviewID, review2.ReviewID)
		assert.Equal(t, *reviewResponse.FromUserID, *review2.FromUserID)
		assert.Equal(t, *reviewResponse.ToUserID, *review2.ToUserID)
		assert.Equal(t, *reviewResponse.Text, *review2.Text)

		// test request with query
		req, errRequest = http.NewRequestWithContext(
			context,
			http.MethodGet,
			fmt.Sprint(testServer.URL, "/api/reviews?to_user_id=", user3.UserID),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse = http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder = json.NewDecoder(res.Body)

		errDecode2 := decoder.Decode(&reviewsFound)
		assert.Nil(t, errDecode2)

		assert.Equal(t, len(reviewsFound), 1)

		reviewResponse = reviewsFound[0]

		assert.Equal(t, reviewResponse.ReviewID, review2.ReviewID)
		assert.Equal(t, *reviewResponse.FromUserID, *review2.FromUserID)
		assert.Equal(t, *reviewResponse.ToUserID, *review2.ToUserID)
		assert.Equal(t, *reviewResponse.Text, *review2.Text)
	})

	t.Run("Test Create", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		user1 := models.User{}
		user2 := models.User{}

		db.Create(&user1)
		db.Create(&user2)

		jsonStr := []byte(fmt.Sprintf(
			`{"text":"Text", "from_user_id": "%s", "to_user_id": "%s"}`,
			&user1.UserID,
			&user2.UserID,
		))

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodPost,
			fmt.Sprint(testServer.URL, "/api/reviews"),
			bytes.NewBuffer(jsonStr),
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var reviewResponse models.Review
		errDecode := decoder.Decode(&reviewResponse)
		assert.Nil(t, errDecode)

		assert.NotNil(t, reviewResponse.ReviewID)
		assert.Equal(t, user1.UserID, *reviewResponse.FromUserID)
		assert.Equal(t, user2.UserID, *reviewResponse.ToUserID)
		assert.Equal(t, "Text", *reviewResponse.Text)

		var reviewFound models.Review
		errFound := db.Where("review_id = ?", reviewResponse.ReviewID).First(&reviewFound).Error

		assert.Nil(t, errFound)

		assert.Equal(t, user1.UserID, *reviewFound.FromUserID)
		assert.Equal(t, user2.UserID, *reviewFound.ToUserID)
		assert.Equal(t, "Text", *reviewFound.Text)
	})

	t.Run("Test Modify", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		user1 := models.User{}
		user2 := models.User{}

		db.Create(&user1)
		db.Create(&user2)

		text := "Text"
		review := models.Review{FromUserID: &user1.UserID, ToUserID: &user2.UserID, Text: &text}

		db.Create(&review)

		jsonStr := []byte(`{"text":"TextDifferent"}`)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodPut,
			fmt.Sprint(testServer.URL, "/api/reviews/", review.ReviewID),
			bytes.NewBuffer(jsonStr),
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var reviewResponse models.Review
		errDecode := decoder.Decode(&reviewResponse)
		assert.Nil(t, errDecode)

		assert.Equal(t, reviewResponse.ReviewID, review.ReviewID)
		assert.Equal(t, "TextDifferent", *reviewResponse.Text)

		var reviewFound models.Review
		errFound := db.Where("review_id = ?", review.ReviewID).First(&reviewFound).Error

		assert.Nil(t, errFound)

		assert.Equal(t, reviewFound.ReviewID, review.ReviewID)
		assert.Equal(t, "TextDifferent", *reviewFound.Text)
	})

	t.Run("Test Delete", func(t *testing.T) {
		errTruncate := dbUtility.TruncateAll()
		assert.Nil(t, errTruncate)

		user1 := models.User{}
		user2 := models.User{}

		db.Create(&user1)
		db.Create(&user2)

		text := "Text"
		review := models.Review{FromUserID: &user1.UserID, ToUserID: &user2.UserID, Text: &text}

		db.Create(&review)

		req, errRequest := http.NewRequestWithContext(
			context,
			http.MethodDelete,
			fmt.Sprint(testServer.URL, "/api/reviews/", review.ReviewID),
			nil,
		)
		assert.Nil(t, errRequest)

		res, errResponse := http.DefaultClient.Do(req)

		assert.Nil(t, errResponse)

		defer res.Body.Close()

		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		var reviewFound models.Review
		errFound := db.Where("review_id = ?", review.ReviewID).First(&reviewFound).Error

		assert.Equal(t, gorm.ErrRecordNotFound, errFound)
	})
}
