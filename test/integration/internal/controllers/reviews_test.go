package reviews_test

import (
	"fmt"
	"context"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"errors"

	"github.com/dgrijalva/jwt-go"

	"gorm.io/gorm"

	"go_server/internal/routes"
	"go_server/internal/models"
	"go_server/internal/db"
	"go_server/internal/middlewares"
)

func init() {
	middlewares.Auth = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwtToken := &jwt.Token{ Claims: jwt.MapClaims{ "sub": "auth0|loggedInReview" } }
	    newContext := context.WithValue(r.Context(), "review", jwtToken)
      h.ServeHTTP(w, r.WithContext(newContext))
	  })
	}
}

func TestGet(t *testing.T) {
	db.Connect()
	db.DB.Exec("truncate table users")
	db.DB.Exec("truncate table reviews")

	testServer := httptest.NewServer(routes.Init())
	defer testServer.Close()

	user1 := models.User{}
	user2 := models.User{}

	db.DB.Create(&user1)
	db.DB.Create(&user2)

	text := "Text"
  review := models.Review{ FromUserId: user1.UserID, ToUserId: user2.UserID, Text: &text }

  db.DB.Create(&review)

	res, errRequest := http.Get(fmt.Sprint(testServer.URL, "/api/reviews/", review.ReviewID))

	if errRequest != nil {
		t.Fatalf("Get: %v", errRequest)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d, Received %d", http.StatusOK, res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)

	var reviewResponse models.Review
	errDecode := decoder.Decode(&reviewResponse)

	if errDecode != nil {
		t.Fatalf("Decoding error: %v", errDecode)
	}

	if review.ReviewID != reviewResponse.ReviewID {
		t.Fatalf("Expected: %s, Received: %s", review.ReviewID, reviewResponse.ReviewID)
	}

	if *review.Text != *reviewResponse.Text {
		t.Fatalf("Expected: %s, Received: %s", *review.Text, *reviewResponse.Text)
	}

	if *review.FromUserID != *reviewResponse.FromUserID {
		t.Fatalf("Expected: %s, Received: %s", *review.FromUserID, *reviewResponse.FromUserID)
	}

	if *review.ToUserID != *reviewResponse.ToUserID {
		t.Fatalf("Expected: %s, Received: %s", *review.ToUserID, *reviewResponse.ToUserID)
	}
}

func TestList(t *testing.T) {
	db.Connect()
	db.DB.Exec("truncate table users")
	db.DB.Exec("truncate table reviews")

	testServer := httptest.NewServer(routes.Init())
	defer testServer.Close()

	user1 := models.User{}
	user2 := models.User{}
	user3 := models.User{}
	user4 := models.User{}

	db.DB.Create(&user1)
	db.DB.Create(&user2)
	db.DB.Create(&user3)
	db.DB.Create(&user4)

	text1 := "Text1"
  review1 := models.Review{ FromUserId: user1.UserID, ToUserId: user2.UserID, Text: &text1 }

	text2 := "Text2"
  review2 := models.Review{ FromUserId: user2.UserID, ToUserId: user3.UserID, Text: &text2 }

	res, errRequest := http.Get(fmt.Sprint(testServer.URL, "/api/reviews"))

	if errRequest != nil {
		t.Fatalf("Get: %v", errRequest)
	}

	if errRequestQuery != nil {
		t.Fatalf("Get: %v", errRequest)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d, Received %d", http.StatusOK, res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)

	var reviewsFound []models.Review
	errDecode := decoder.Decode(&reviewsFound)

	if errDecode != nil {
		t.Fatalf("Decoding error: %v", errDecode)
	}

	if len(reviewsFound) != 2 {
		t.Fatalf("Expected: %d, Received: %d", 2, len(reviewsFound))
	}

	var reviewResponse models.Review

	for _, value := range reviewsFound {
    if value.ReviewID == review1.ReviewID {
			reviewResponse = value
			break
    }
	}

	if review1.ReviewID != reviewResponse.ReviewID {
		t.Fatalf("Expected: %s, Received: %s", review1.ReviewID, *reviewResponse.ReviewID)
	}

	if *review1.FromUserId != *reviewResponse.FromUserId {
		t.Fatalf("Expected: %s, Received: %s", *review1.FromUserId, *reviewResponse.FromUserId)
	}

	if *review1.ToUserId != *reviewResponse.ToUserId {
		t.Fatalf("Expected: %s, Received: %s", *review1.ToUserId, *reviewResponse.ToUserId)
	}

	if *review1.Text != *reviewResponse.Text {
		t.Fatalf("Expected: %s, Received: %s", *review1.Text, *reviewResponse.Text)
	}

	for _, value := range reviewsFound {
    if value.ReviewID == review2.ReviewID {
			reviewResponse = value
			break
    }
	}

	if review2.ReviewID != reviewResponse.ReviewID {
		t.Fatalf("Expected: %s, Received: %s", review2.ReviewID, *reviewResponse.ReviewID)
	}

	if *review2.FromUserId != *reviewResponse.FromUserId {
		t.Fatalf("Expected: %s, Received: %s", *review2.FromUserId, *reviewResponse.FromUserId)
	}

	if *review2.ToUserId != *reviewResponse.ToUserId {
		t.Fatalf("Expected: %s, Received: %s", *review2.ToUserId, *reviewResponse.ToUserId)
	}

	if *review2.Text != *reviewResponse.Text {
		t.Fatalf("Expected: %s, Received: %s", *review2.Text, *reviewResponse.Text)
	}

	// test request with query

	res, errRequest := http.Get(fmt.Sprint(testServer.URL, "/api/reviews?to_user_id=", user3.UserID))

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d, Received %d", http.StatusOK, res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)

	var reviewsFound []models.Review
	errDecode := decoder.Decode(&reviewsFound)

	if errDecode != nil {
		t.Fatalf("Decoding error: %v", errDecode)
	}

	if len(reviewsFound) != 1 {
		t.Fatalf("Expected: %d, Received: %d", 1, len(reviewsFound))
	}

	reivewResponse = reviewsFound[0]

	if review2.ReviewID != reviewResponse.ReviewID {
		t.Fatalf("Expected: %s, Received: %s", review2.ReviewID, *reviewResponse.ReviewID)
	}

	if *review2.FromUserId != *reviewResponse.FromUserId {
		t.Fatalf("Expected: %s, Received: %s", *review2.FromUserId, *reviewResponse.FromUserId)
	}

	if *review2.ToUserId != *reviewResponse.ToUserId {
		t.Fatalf("Expected: %s, Received: %s", *review2.ToUserId, *reviewResponse.ToUserId)
	}

	if *review2.Text != *reviewResponse.Text {
		t.Fatalf("Expected: %s, Received: %s", *review2.Text, *reviewResponse.Text)
	}
}

func TestCreate(t *testing.T) {
	db.Connect()
	db.DB.Exec("truncate table users")
	db.DB.Exec("truncate table reviews")

	testServer := httptest.NewServer(routes.Init())
	defer testServer.Close()

	user1 := models.User{}
	user2 := models.User{}

	db.DB.Create(&user1)
	db.DB.Create(&user2)

	var jsonStr = []byte(fmt.Sprint(`{"text":"Text", "from_user_id": "%s", "to_user_id": "%s"}`, &user1.UserID, &user2.UserID))

	res, errRequest := http.Post(fmt.Sprint(testServer.URL, "/api/reviews"), "application/json", bytes.NewBuffer(jsonStr))

	if errRequest != nil {
		t.Fatalf("Get: %v", errRequest)
	}

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Expected %d, Received %d", http.StatusCreated, res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)

	var reviewResponse models.Review
	errDecode := decoder.Decode(&reviewResponse)

	if errDecode != nil {
		t.Fatalf("Decoding error: %v", errDecode)
	}

	if *reviewResponse.Text != "Text" {
		t.Fatalf("Expected: %s, Received: %s", "Text", *reviewResponse.Text)
	}

	if *reviewResponse.FromUserID != user1.UserId {
		t.Fatalf("Expected: %s, Received: %s", user1.UserId, *reviewResponse.FromUserID)
	}

	if *reviewResponse.ToUserID != user2.UserId {
		t.Fatalf("Expected: %s, Received: %s", user2.UserId, *reviewResponse.ToUserID)
	}

	var reviewFound models.Review
	errFound := db.DB.Where("review_id = ?", reviewResponse.ReviewID).First(&reviewFound).Error

	if errFound != nil {
		t.Fatalf("Error: %v", errFound)
	}

	if *reviewFound.Text != "Text" {
		t.Fatalf("Expected: %s, Received: %s", "Text", *reviewFound.Text)
	}

	if *reviewFound.FromUserId != user1.UserID {
		t.Fatalf("Expected: %s, Received: %s", user1.UserID, *reviewFound.FromUserId)
	}

	if *reviewFound.ToUserId != user2.UserID {
		t.Fatalf("Expected: %s, Received: %s", user2.UserID, *reviewFound.ToUserId)
	}
}

// func TestModify(t *testing.T) {
// 	db.Connect()
// 	db.DB.Exec("truncate table reviews")
// 
// 	testServer := httptest.NewServer(routes.Init())
// 	defer testServer.Close()
//
// 	firstName:= "firstName"
//   review := models.Review{ FirstName: &firstName }
//
//   db.DB.Create(&review)
//
//
// 	firstNameDifferent := "firstNameDifferent"
// 	reviewDifferent := models.Review{ FirstName: &firstNameDifferent }
// 	jsonReq, errJSON := json.Marshal(reviewDifferent)
//
// 	if errJSON != nil {
// 		t.Fatalf("JSON: %v", errJSON)
// 	}
//
// 	req, errRequest := http.NewRequest(http.MethodPut, fmt.Sprint(testServer.URL, "/api/reviews/", review.ReviewID), bytes.NewBuffer(jsonReq))
// 	client := &http.Client{}
// 	res, errRequest := client.Do(req)
//
// 	if errRequest != nil {
// 		t.Fatalf("Delete: %v", errRequest)
// 	}
//
// 	if res.StatusCode != http.StatusOK {
//     t.Fatalf("Expected %d, Received %d", http.StatusOK, res.StatusCode)
//   }
//
// 	decoder := json.NewDecoder(res.Body)
//
// 	var reviewResponse models.Review
// 	errDecode := decoder.Decode(&reviewResponse)
//
// 	if errDecode != nil {
// 		t.Fatalf("Decoding error: %v", errDecode)
// 	}
//
// 	if review.ReviewID != reviewResponse.ReviewID {
// 		t.Fatalf("Expected: %s, Received: %s", review.ReviewID, reviewResponse.ReviewID)
// 	}
//
// 	if *reviewDifferent.FirstName != *reviewResponse.FirstName {
// 		t.Fatalf("Expected: %s, Received: %s", *review.FirstName, *reviewResponse.FirstName)
// 	}
//
// 	var reviewFound models.Review
//   errFound := db.DB.Where("review_id = ?", review.ReviewID).First(&reviewFound).Error
//
// 	if errFound != nil {
// 		t.Fatalf("Error: %v", errFound)
// 	}
//
// 	if review.ReviewID != reviewFound.ReviewID {
// 		t.Fatalf("Expected: %s, Received: %s", review.ReviewID, reviewFound.ReviewID)
// 	}
//
// 	if *reviewDifferent.FirstName != *reviewFound.FirstName {
// 		t.Fatalf("Expected: %s, Received: %s", *review.FirstName, *reviewFound.FirstName)
// 	}
// }
//
// func TestDelete(t *testing.T) {
// 	db.Connect()
// 	db.DB.Exec("truncate table reviews")
//
// 	testServer := httptest.NewServer(routes.Init())
// 	defer testServer.Close()
//
// 	firstName := "firstName"
//   review := models.Review{ FirstName: &firstName }
//
//   db.DB.Create(&review)
//
// 	req, errRequest := http.NewRequest(http.MethodDelete, fmt.Sprint(testServer.URL, "/api/reviews/", review.ReviewID), nil)
// 	client := &http.Client{}
// 	res, errRequest := client.Do(req)
//
// 	if errRequest != nil {
// 		t.Fatalf("Delete: %v", errRequest)
// 	}
//
// 	if res.StatusCode != http.StatusNoContent {
//     t.Fatalf("Expected %d, Received %d", http.StatusNoContent, res.StatusCode)
//   }
//
// 	var reviewFound models.Review
//   errFound := db.DB.Where("review_id = ?", review.ReviewID).First(&reviewFound).Error
//
// 	if !errors.Is(errFound, gorm.ErrRecordNotFound) {
// 		t.Fatalf("Expected review not to be found")
// 	}
// }
