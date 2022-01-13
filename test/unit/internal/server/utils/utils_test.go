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

  "github.com/stretchr/testify/assert"
  "github.com/google/uuid"
)

func TestUtilsGetQueryParamUUID(t *testing.T) {
	t.Parallel()

	zapLogger := zaptest.NewLogger(t)
	logger := goServerZapLogger.NewLogger(zapLogger)

  utils := utils.NewServerUtils(logger)

  uuid := uuid.New()
  req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://www.google.com?testParam=%s", uuid), nil)
  uuidFound := utils.GetQueryParamUUID(req, "testParam")

  assert.Equal(t, uuid, uuidFound)
}

func TestUtilsGetQueryParamUUIDMe(t *testing.T) {
	t.Parallel()

	zapLogger := zaptest.NewLogger(t)
	logger := goServerZapLogger.NewLogger(zapLogger)

  utils := utils.NewServerUtils(logger)

  uuid := uuid.New()
  user := models.User{ UserID: uuid }
  req := httptest.NewRequest(http.MethodGet, "http://www.google.com?testParam=me", nil)
  newContext := context.WithValue(req.Context(), consts.UserModelKey, user)
  req = req.WithContext(newContext)

  uuidFound := utils.GetQueryParamUUID(req, "testParam")

  assert.Equal(t, uuid, uuidFound)
}

func TestUtilsGetQueryParamUUIDUndifined(t *testing.T) {
	t.Parallel()

	zapLogger := zaptest.NewLogger(t)
	logger := goServerZapLogger.NewLogger(zapLogger)

  utils := utils.NewServerUtils(logger)

  req := httptest.NewRequest(http.MethodGet, "http://www.google.com", nil)

  uuidFound := utils.GetQueryParamUUID(req, "testParam")

  assert.Equal(t, uuid.Nil, uuidFound)
}
