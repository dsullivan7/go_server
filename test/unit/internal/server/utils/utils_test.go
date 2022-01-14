package utils_test

import (
  "context"
	"fmt"
  "net/http"
  "net/http/httptest"
  "testing"

  goServerZapLogger "go_server/internal/logger/zap"
	"go_server/internal/server/utils"
  "go.uber.org/zap/zaptest"
  "go_server/internal/models"
  "go_server/internal/server/consts"

  "github.com/go-chi/chi"
  "github.com/stretchr/testify/assert"
  "github.com/google/uuid"
)

type ctxKey struct {
  name string
}

func TestUtils(t *testing.T) {
  t.Parallel()

  zapLogger := zaptest.NewLogger(t)
	logger := goServerZapLogger.NewLogger(zapLogger)

  utils := utils.NewServerUtils(logger)

  t.Run("GetPathParamUUID", func(t *testing.T) {
    uuid := uuid.New()
    req := httptest.NewRequest(http.MethodGet, "https://google.com", nil)
    rctx := chi.NewRouteContext()
    rctx.URLParams.Add("testParam", uuid.String())
    newContext := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
    req = req.WithContext(newContext)
    uuidFound := utils.GetPathParamUUID(req, "testParam")

    assert.Equal(t, uuid, uuidFound)
  })

  t.Run("GetPathParamUUID me", func(t *testing.T) {
    uuid := uuid.New()
    user := models.User{ UserID: uuid }
    req := httptest.NewRequest(http.MethodGet, "https://google.com", nil)
    newContext := context.WithValue(req.Context(), consts.UserModelKey, user)
    req = req.WithContext(newContext)
    rctx := chi.NewRouteContext()
    rctx.URLParams.Add("testParam", uuid.String())
    newContext = context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
    req = req.WithContext(newContext)
    uuidFound := utils.GetPathParamUUID(req, "testParam")

    assert.Equal(t, uuid, uuidFound)
  })

  t.Run("GetQueryParamUUID", func(t *testing.T) {
    uuid := uuid.New()
    req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("https://google.com?testParam=%s", uuid), nil)
    uuidFound := utils.GetQueryParamUUID(req, "testParam")

    assert.Equal(t, uuid, uuidFound)
  })

  t.Run("GetQueryParamUUID me", func(t *testing.T) {
    uuid := uuid.New()
    user := models.User{ UserID: uuid }
    req := httptest.NewRequest(http.MethodGet, "https://google.com?testParam=me", nil)
    newContext := context.WithValue(req.Context(), consts.UserModelKey, user)
    req = req.WithContext(newContext)

    uuidFound := utils.GetQueryParamUUID(req, "testParam")

    assert.Equal(t, uuid, uuidFound)
  })

  t.Run("GetQueryParamUUID undefined", func(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "https://google.com", nil)

    uuidFound := utils.GetQueryParamUUID(req, "testParam")

    assert.Equal(t, uuid.Nil, uuidFound)
  })
}
