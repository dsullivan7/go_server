package middlewares_test

import (
  "context"
  "net/http"
  "net/http/httptest"
  "testing"

  goServerZapLogger "go_server/internal/logger/zap"
	"go_server/internal/server/utils"
  "go.uber.org/zap/zaptest"
  "go_server/internal/models"
  "go_server/internal/server/consts"
  "go_server/internal/server/middlewares"
  "go_server/internal/config"
  "go_server/test/mocks/store"
  "go_server/test/mocks/auth"

  "github.com/go-chi/chi"
  "github.com/stretchr/testify/assert"
  "github.com/google/uuid"
)

const domain = "https://sunburst.app"

func TestMiddlewares(t *testing.T) {
  t.Parallel()

  cfg, configError := config.NewConfig()
  assert.Nil(t, configError)

  zapLogger := zaptest.NewLogger(t)
	logger := goServerZapLogger.NewLogger(zapLogger)

  utils := utils.NewServerUtils(logger)

  storeMock := store.NewMockStore()

  authMock := auth.NewMockAuth()

  mw := middlewares.NewMiddlewares(cfg, storeMock, authMock, utils, logger)

  t.Run("PathParam", func(t *testing.T) {
    uuid := uuid.New()
    user := models.User{ UserID: uuid }

    param := "testParam"

    nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      paramValue := chi.URLParam(r, param)
      assert.Equal(t, uuid.String(), paramValue)
    })

    paramMiddleware := mw.URLParam("testParam")
    testHandler := paramMiddleware(nextHandler)

    req := httptest.NewRequest(http.MethodGet, domain, nil)
    newContext := context.WithValue(req.Context(), consts.UserModelKey, user)

    rctx := chi.NewRouteContext()
    rctx.URLParams.Add(param, "me")
    newContext = context.WithValue(newContext, chi.RouteCtxKey, rctx)
    req = req.WithContext(newContext)

    testHandler.ServeHTTP(nil, req)
  })
}
