package integration_test

import (
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
	"gorm.io/gorm"
)

func TestReviews(t *testing.T) {
	config, configError := config.NewConfig()
	assert.Nil(t, configError)

	logger, errLogger := logger.NewZapLogger()
	assert.Nil(t, errLogger)

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
		db.Exec("truncate table reviews cascade")
		db.Exec("truncate table users cascade")

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

		assert.Equal(t, res.StatusCode, http.StatusOK)

		decoder := json.NewDecoder(res.Body)

		var reviewResponse models.Review

		errDecode := decoder.Decode(&reviewResponse)
		assert.Nil(t, errDecode)

		assert.Equal(t, review.ReviewID, reviewResponse.ReviewID)
		assert.Equal(t, *review.Text, *reviewResponse.Text)
		assert.Equal(t, *review.FromUserID, *reviewResponse.FromUserID)
		assert.Equal(t, *review.ToUserID, *reviewResponse.ToUserID)
	})

	t.Run("Test List", func(t *testing.T) {
		db.Exec("truncate table reviews cascade")
		db.Exec("truncate table users cascade")

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

		assert.Equal(t, res.StatusCode, http.StatusOK)

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

		assert.Equal(t, review1.ReviewID, reviewResponse.ReviewID)
		assert.Equal(t, *review1.FromUserID, *reviewResponse.FromUserID)
		assert.Equal(t, *review1.ToUserID, *reviewResponse.ToUserID)
		assert.Equal(t, *review1.Text, *reviewResponse.Text)

		for _, value := range reviewsFound {
			if value.ReviewID == review2.ReviewID {
				reviewResponse = value

				break
			}
		}

		assert.Equal(t, review2.ReviewID, reviewResponse.ReviewID)
		assert.Equal(t, *review2.FromUserID, *reviewResponse.FromUserID)
		assert.Equal(t, *review2.ToUserID, *reviewResponse.ToUserID)
		assert.Equal(t, *review2.Text, *reviewResponse.Text)

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

		assert.Equal(t, res.StatusCode, http.StatusOK)

		decoder = json.NewDecoder(res.Body)

		errDecode2 := decoder.Decode(&reviewsFound)
		assert.Nil(t, errDecode2)

		assert.Equal(t, len(reviewsFound), 1)

		reviewResponse = reviewsFound[0]

		assert.Equal(t, review2.ReviewID, reviewResponse.ReviewID)
		assert.Equal(t, *review2.FromUserID, *reviewResponse.FromUserID)
		assert.Equal(t, *review2.ToUserID, *reviewResponse.ToUserID)
		assert.Equal(t, *review2.Text, *reviewResponse.Text)
	})

	t.Run("Test Create", func(t *testing.T) {
		db.Exec("truncate table reviews cascade")
		db.Exec("truncate table users cascade")

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

		assert.Equal(t, res.StatusCode, http.StatusCreated)

		decoder := json.NewDecoder(res.Body)

		var reviewResponse models.Review
		errDecode := decoder.Decode(&reviewResponse)
		assert.Nil(t, errDecode)

		assert.NotNil(t, reviewResponse.ReviewID)
		assert.Equal(t, *reviewResponse.FromUserID, user1.UserID)
		assert.Equal(t, *reviewResponse.ToUserID, user2.UserID)
		assert.Equal(t, *reviewResponse.Text, "Text")

		var reviewFound models.Review
		errFound := db.Where("review_id = ?", reviewResponse.ReviewID).First(&reviewFound).Error

		assert.Nil(t, errFound)

		assert.Equal(t, *reviewFound.FromUserID, user1.UserID)
		assert.Equal(t, *reviewFound.ToUserID, user2.UserID)
		assert.Equal(t, *reviewFound.Text, "Text")
	})

	t.Run("Test Modify", func(t *testing.T) {
		db.Exec("truncate table reviews cascade")
		db.Exec("truncate table users cascade")

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

		assert.Equal(t, res.StatusCode, http.StatusOK)

		decoder := json.NewDecoder(res.Body)

		var reviewResponse models.Review
		errDecode := decoder.Decode(&reviewResponse)
		assert.Nil(t, errDecode)

		assert.Equal(t, review.ReviewID, reviewResponse.ReviewID)
		assert.Equal(t, *reviewResponse.Text, "TextDifferent")

		var reviewFound models.Review
		errFound := db.Where("review_id = ?", review.ReviewID).First(&reviewFound).Error

		assert.Nil(t, errFound)

		assert.Equal(t, review.ReviewID, reviewFound.ReviewID)
		assert.Equal(t, *reviewFound.Text, "TextDifferent")
	})

	t.Run("Test Delete", func(t *testing.T) {
		db.Exec("truncate table reviews cascade")
		db.Exec("truncate table users cascade")

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

		assert.Equal(t, res.StatusCode, http.StatusNoContent)

		var reviewFound models.Review
		errFound := db.Where("review_id = ?", review.ReviewID).First(&reviewFound).Error

		assert.Equal(t, errFound, gorm.ErrRecordNotFound)
	})
}
