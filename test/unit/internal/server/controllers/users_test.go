package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"go_server/internal/models"
	"go_server/internal/server/consts"
	"go_server/test/unit/internal/server/controllers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestUsers(tParent *testing.T) {
	tParent.Parallel()

	controllers, mockStore, err := controllers.Setup()
	assert.Nil(tParent, err)

	tParent.Run("Test Get", func(t *testing.T) {
		t.Parallel()

		firstName := "firstName"
		lastName := "lastName"
		auth0ID := "auth0ID"

		user := models.User{
			FirstName: &firstName,
			LastName:  &lastName,
			Auth0ID:   &auth0ID,
		}

		uuid := uuid.New()

		mockStore.On("GetUser", uuid).Return(&user, nil)

		req := httptest.NewRequest(
			http.MethodGet,
			"/api/users",
			nil,
		)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("user_id", uuid.String())

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		controllers.GetUser(w, req)

		res := w.Result()
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

	tParent.Run("Test Get Me", func(t *testing.T) {
		t.Parallel()

		firstName := "firstName"
		lastName := "lastName"
		auth0ID := "auth0ID"

		user := models.User{
			UserID:    uuid.New(),
			FirstName: &firstName,
			LastName:  &lastName,
			Auth0ID:   &auth0ID,
		}

		mockStore.On("GetUser", user.UserID).Return(&user, nil)

		req := httptest.NewRequest(
			http.MethodGet,
			"/api/users",
			nil,
		)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("user_id", "me")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		req = req.WithContext(context.WithValue(req.Context(), consts.UserModelKey, user))

		w := httptest.NewRecorder()

		controllers.GetUser(w, req)

		res := w.Result()
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

	tParent.Run("Test List", func(t *testing.T) {
		t.Parallel()

		firstName1 := "firstName1"
		lastName1 := "lastName1"
		auth0Id1 := "auth0Id1"

		firstName2 := "firstName2"
		lastName2 := "lastName2"
		auth0Id2 := "auth0Id2"

		user1 := models.User{
			UserID:    uuid.New(),
			FirstName: &firstName1,
			LastName:  &lastName1,
			Auth0ID:   &auth0Id1,
		}

		user2 := models.User{
			UserID:    uuid.New(),
			FirstName: &firstName2,
			LastName:  &lastName2,
			Auth0ID:   &auth0Id2,
		}

		mockStore.On("ListUsers", map[string]interface{}{}).Return([]models.User{user1, user2}, nil)

		req := httptest.NewRequest(
			http.MethodGet,
			"/api/users",
			nil,
		)

		w := httptest.NewRecorder()

		controllers.ListUsers(w, req)

		res := w.Result()
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

	tParent.Run("Test Create", func(t *testing.T) {
		t.Parallel()

		jsonStr := []byte(`{
			"first_name":"firstName",
			"last_name":"lastName",
			"auth0_id":"auth0Id"
		}`)

		firstName := "firstName"
		lastName := "lastName"
		auth0Id := "auth0Id"

		userPayload := models.User{
			FirstName: &firstName,
			LastName:  &lastName,
			Auth0ID:   &auth0Id,
		}

		userCreated := models.User{
			UserID:    uuid.New(),
			FirstName: &firstName,
			LastName:  &lastName,
			Auth0ID:   &auth0Id,
		}

		mockStore.On("CreateUser", userPayload).Return(&userCreated, nil)

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/users",
			bytes.NewBuffer(jsonStr),
		)

		w := httptest.NewRecorder()

		controllers.CreateUser(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

		decoder := json.NewDecoder(res.Body)

		var userResponse models.User
		errDecoder := decoder.Decode(&userResponse)
		assert.Nil(t, errDecoder)

		assert.NotNil(t, userResponse.UserID)
		assert.Equal(t, "firstName", *userResponse.FirstName)
		assert.Equal(t, "lastName", *userResponse.LastName)
		assert.Equal(t, "auth0Id", *userResponse.Auth0ID)
	})

	// tParent.Run("Test Modify", func(t *testing.T) {
	// 	t.Parallel()
	//
	// 	firstName := "FirstName"
	// 	lastName := "LastName"
	// 	auth0ID := "Auth0ID"
	// 	user := models.User{
	// 		FirstName: &firstName,
	// 		LastName:  &lastName,
	// 		Auth0ID:   &auth0ID,
	// 	}
	//
	// 	db.Create(&user)
	//
	// 	jsonStr := []byte(`{
	// 		"first_name":"FirstNameDifferent",
	// 		"last_name": "LastNameDifferent",
	// 		"auth0_id": "Auth0IDDifferent"
	// 	}`)
	//
	// 	req, errRequest := http.NewRequestWithContext(
	// 		context,
	// 		http.MethodPut,
	// 		fmt.Sprint(testServer.URL, "/api/users/", user.UserID),
	// 		bytes.NewBuffer(jsonStr),
	// 	)
	// 	assert.Nil(t, errRequest)
	//
	// 	res, errResponse := http.DefaultClient.Do(req)
	//
	// 	assert.Nil(t, errResponse)
	//
	// 	defer res.Body.Close()
	//
	// 	assert.Equal(t, http.StatusOK, res.StatusCode)
	// 	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))
	//
	// 	decoder := json.NewDecoder(res.Body)
	//
	// 	var userResponse models.User
	// 	errDecoder := decoder.Decode(&userResponse)
	// 	assert.Nil(t, errDecoder)
	//
	// 	assert.Equal(t, "FirstNameDifferent", *userResponse.FirstName)
	// 	assert.Equal(t, "LastNameDifferent", *userResponse.LastName)
	// 	assert.Equal(t, "Auth0IDDifferent", *userResponse.Auth0ID)
	//
	// 	var userFound models.User
	// 	errFound := db.Where("user_id = ?", user.UserID).First(&userFound).Error
	//
	// 	assert.Nil(t, errFound)
	//
	// 	assert.Equal(t, "FirstNameDifferent", *userFound.FirstName)
	// 	assert.Equal(t, "LastNameDifferent", *userFound.LastName)
	// 	assert.Equal(t, "Auth0IDDifferent", *userFound.Auth0ID)
	// })
	//
	// tParent.Run("Test Delete", func(t *testing.T) {
	// 	t.Parallel()
	//
	// 	firstName := "firstName"
	// 	user := models.User{FirstName: &firstName}
	//
	// 	db.Create(&user)
	//
	// 	req, errRequest := http.NewRequestWithContext(
	// 		context,
	// 		http.MethodDelete,
	// 		fmt.Sprint(testServer.URL, "/api/users/", user.UserID),
	// 		nil,
	// 	)
	// 	assert.Nil(t, errRequest)
	//
	// 	res, errResponse := http.DefaultClient.Do(req)
	//
	// 	assert.Nil(t, errResponse)
	//
	// 	defer res.Body.Close()
	//
	// 	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	//
	// 	var userFound models.User
	// 	errFound := db.Where("user_id = ?", user.UserID).First(&userFound).Error
	//
	// 	assert.Equal(t, errFound, gorm.ErrRecordNotFound)
	// })
}
