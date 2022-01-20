package controllers_test

import (
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

	ctx := context.Background()

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

		req, errRequest := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"/api/users/",
			nil,
		)
		assert.Nil(t, errRequest)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("user_id", uuid.String())

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		res := httptest.NewRecorder()

		handler := http.HandlerFunc(controllers.GetUser)
		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

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
      UserID: uuid.New(),
			FirstName: &firstName,
			LastName:  &lastName,
			Auth0ID:   &auth0ID,
		}

		mockStore.On("GetUser", user.UserID).Return(&user, nil)

		req, errRequest := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"/api/users/",
			nil,
		)
		assert.Nil(t, errRequest)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("user_id", "me")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
    req = req.WithContext(context.WithValue(req.Context(), consts.UserModelKey, user))

		res := httptest.NewRecorder()

		handler := http.HandlerFunc(controllers.GetUser)
		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		decoder := json.NewDecoder(res.Body)

		var userResponse models.User
		errDecoder := decoder.Decode(&userResponse)
		assert.Nil(t, errDecoder)

		assert.Equal(t, userResponse.UserID, user.UserID)
		assert.Equal(t, *userResponse.FirstName, *user.FirstName)
		assert.Equal(t, *userResponse.LastName, *user.LastName)
		assert.Equal(t, *userResponse.Auth0ID, *user.Auth0ID)
	})
}
