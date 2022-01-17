package integration_test

import (
	"context"
	"encoding/json"
	"go_server/internal/captcha/twocaptcha"
	"go_server/internal/config"
	goServerRodCrawler "go_server/internal/crawler/rod"
	goServerZapLogger "go_server/internal/logger/zap"
	"go_server/internal/models"
	"go_server/internal/server/controllers"
	"go_server/internal/server/utils"
	"go_server/test/mocks/store"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/go-chi/chi"
	"github.com/go-rod/rod"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestUsers(tParent *testing.T) {
	tParent.Parallel()

	config, configError := config.NewConfig()
	assert.Nil(tParent, configError)

	zapLogger, errZap := zap.NewProduction()
	assert.Nil(tParent, errZap)

	logger := goServerZapLogger.NewLogger(zapLogger)

	store := store.NewMockStore()

	browser := rod.New()

	captchaKey := "key"

	captcha := twocaptcha.NewTwoCaptcha(captchaKey, logger)

	crawler := goServerRodCrawler.NewCrawler(browser, captcha)

	ctx := context.Background()

	utils := utils.NewServerUtils(logger)
	controllers := controllers.NewControllers(config, store, crawler, utils, logger)

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

		store.On("GetUser", uuid).Return(&user, nil)

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
}
