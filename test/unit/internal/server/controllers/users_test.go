package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go_server/internal/models"
	"go_server/internal/server/consts"
	"go_server/test/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestUserGet(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	firstName := "firstName"
	lastName := "lastName"
	auth0ID := "auth0ID"

	user := models.User{
		FirstName: &firstName,
		LastName:  &lastName,
		Auth0ID:   &auth0ID,
	}

	uuid := uuid.New()

	testServer.Store.On("GetUser", uuid).Return(&user, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/users",
		nil,
	)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("user_id", uuid.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().GetUser(w, req)

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

	testServer.Store.AssertExpectations(t)
}

func TestUserGetMe(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	firstName := "firstName"
	lastName := "lastName"
	auth0ID := "auth0ID"

	user := models.User{
		UserID:    uuid.New(),
		FirstName: &firstName,
		LastName:  &lastName,
		Auth0ID:   &auth0ID,
	}

	testServer.Store.On("GetUser", user.UserID).Return(&user, nil)

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

	testServer.Server.GetControllers().GetUser(w, req)

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

	testServer.Store.AssertExpectations(t)
}

func TestUserList(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

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

	testServer.Store.On("ListUsers", map[string]interface{}{}).Return([]models.User{user1, user2}, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/users",
		nil,
	)

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().ListUsers(w, req)

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

	testServer.Store.AssertExpectations(t)
}

func TestUserCreate(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	firstName := "firstName"
	lastName := "lastName"
	auth0Id := "auth0Id"

	jsonStr := []byte(fmt.Sprintf(
		`{
				"first_name":"%s",
				"last_name":"%s",
				"auth0_id":"%s"
			}`,
		firstName,
		lastName,
		auth0Id,
	))

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

	testServer.Store.On("CreateUser", userPayload).Return(&userCreated, nil)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/users",
		bytes.NewBuffer(jsonStr),
	)

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().CreateUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var userResponse models.User
	errDecoder := decoder.Decode(&userResponse)
	assert.Nil(t, errDecoder)

	assert.NotNil(t, userResponse.UserID)
	assert.Equal(t, firstName, *userResponse.FirstName)
	assert.Equal(t, lastName, *userResponse.LastName)
	assert.Equal(t, auth0Id, *userResponse.Auth0ID)

	testServer.Store.AssertExpectations(t)
}

func TestUserModify(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	userID := uuid.New()

	firstNameNew := "firstNameNew"
	lastNameNew := "lastNameNew"
	auth0IdNew := "auth0IdNew"

	jsonStr := []byte(fmt.Sprintf(
		`{
				"first_name":"%s",
				"last_name":"%s",
				"auth0_id":"%s"
			}`,
		firstNameNew,
		lastNameNew,
		auth0IdNew,
	))

	userPayload := models.User{
		FirstName: &firstNameNew,
		LastName:  &lastNameNew,
		Auth0ID:   &auth0IdNew,
	}

	userModified := models.User{
		UserID:    userID,
		FirstName: &firstNameNew,
		LastName:  &lastNameNew,
		Auth0ID:   &auth0IdNew,
	}

	testServer.Store.On("ModifyUser", userID, userPayload).Return(&userModified, nil)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/users",
		bytes.NewBuffer(jsonStr),
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("user_id", userID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().ModifyUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var userResponse models.User
	errDecoder := decoder.Decode(&userResponse)
	assert.Nil(t, errDecoder)

	assert.Equal(t, firstNameNew, *userResponse.FirstName)
	assert.Equal(t, lastNameNew, *userResponse.LastName)
	assert.Equal(t, auth0IdNew, *userResponse.Auth0ID)

	testServer.Store.AssertExpectations(t)
}

func TestUserDelete(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	userID := uuid.New()

	testServer.Store.On("DeleteUser", userID).Return(nil)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/users",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("user_id", userID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().DeleteUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	testServer.Store.AssertExpectations(t)
}
