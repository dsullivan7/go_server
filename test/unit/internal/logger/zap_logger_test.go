package logger_test

import (
	"go_server/internal/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestZapLogger(t *testing.T) {
	t.Parallel()

	zapLogger, errZap := zap.NewProduction()
	assert.Nil(t, errZap)

	logger := logger.NewZapLogger(zapLogger)

	m := make(map[string]interface{})
	m["some"] = "key"

	logger.InfoWithMeta("blah", m)
}
