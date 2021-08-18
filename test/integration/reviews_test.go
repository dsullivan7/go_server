package integration

import (
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"gorm.io/gorm"
	"go_server/internal/config"
	"go_server/internal/models"
	"go_server/internal/store"
	"go_server/internal/controllers"
	"go_server/internal/server"
	"go_server/internal/logger"
	"go_server/internal/db"
	"github.com/go-chi/chi"


	"github.com/stretchr/testify/assert"
)

func TestReviews(t *testing.T) {
	config := config.NewConfig()
	logger := logger.NewZapLogger()

	db, _ := db.NewDatabase(
		config.DBHost,
		config.DBName,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBSSL,
	)

	store := store.NewGormStore(db)

	controllers := controllers.NewControllers(store, config, logger)
	router := chi.NewRouter()
	server := server.NewServer(router, controllers, config, logger)
	testServer := httptest.NewServer(server.Routes())

	defer testServer.Close()

	t.Run("Test Get", func(t *testing.T) {
		db.Exec("truncate table reviews cascade")
		db.Exec("truncate table users cascade")

		user1 := models.User{}
		user2 := models.User{}

		db.Create(&user1)
		db.Create(&user2)

		text := "Text"
	  review := models.Review{ FromUserID: &user1.UserID, ToUserID: &user2.UserID, Text: &text }

	  db.Create(&review)

		res, errRequest := http.Get(fmt.Sprint(testServer.URL, "/api/reviews/", review.ReviewID))

		assert.Nil(t, errRequest)

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
	  review1 := models.Review{ FromUserID: &user1.UserID, ToUserID: &user2.UserID, Text: &text1 }

		text2 := "Text2"
	  review2 := models.Review{ FromUserID: &user2.UserID, ToUserID: &user3.UserID, Text: &text2 }

		db.Create(&review1)
		db.Create(&review2)

		res, errRequest := http.Get(fmt.Sprint(testServer.URL, "/api/reviews"))

		assert.Nil(t, errRequest)

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
		res, errRequest = http.Get(fmt.Sprint(testServer.URL, "/api/reviews?to_user_id=", user3.UserID))

		assert.Nil(t, errRequest)

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

		var jsonStr = []byte(fmt.Sprintf(`{"text":"Text", "from_user_id": "%s", "to_user_id": "%s"}`, &user1.UserID, &user2.UserID))

		res, errRequest := http.Post(fmt.Sprint(testServer.URL, "/api/reviews"), "application/json", bytes.NewBuffer(jsonStr))

		assert.Nil(t, errRequest)

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
	  review := models.Review{ FromUserID: &user1.UserID, ToUserID: &user2.UserID, Text: &text }

		db.Create(&review)

		var jsonStr = []byte(`{"text":"TextDifferent"}`)

		req, errCreate := http.NewRequest(http.MethodPut, fmt.Sprint(testServer.URL, "/api/reviews/", review.ReviewID), bytes.NewBuffer(jsonStr))
		assert.Nil(t, errCreate)

		client := &http.Client{}
		res, errRequest := client.Do(req)

		assert.Nil(t, errRequest)

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
	  review := models.Review{ FromUserID: &user1.UserID, ToUserID: &user2.UserID, Text: &text }

		db.Create(&review)

		req, errCreate := http.NewRequest(http.MethodDelete, fmt.Sprint(testServer.URL, "/api/reviews/", review.ReviewID), nil)
		assert.Nil(t, errCreate)

		client := &http.Client{}
		res, errRequest := client.Do(req)

		assert.Nil(t, errRequest)

		assert.Equal(t, res.StatusCode, http.StatusNoContent)

		var reviewFound models.Review
	  errFound := db.Where("review_id = ?", review.ReviewID).First(&reviewFound).Error

		assert.Equal(t, errFound, gorm.ErrRecordNotFound)
	})
}
