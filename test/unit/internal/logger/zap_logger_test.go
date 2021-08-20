package logger_test

import (
	"go_server/internal/logger"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestZapLogger(t *testing.T) {
	t.Parallel()

	core, recordedLogs := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(core)

	logger := logger.NewZapLogger(zapLogger)

	timeNow := time.Now()
	timeSince := time.Since(timeNow)

	logger.InfoWithMeta(
		"someMessage",
		map[string]interface{}{
			"someString":   "someStringValue",
			"someInt":      52,
			"someDuration": timeSince,
			"someAny":      map[string]interface{}{"someKey": "someValue"},
		},
	)

	logs := recordedLogs.All()

	assert.Equal(t, 1, len(logs))
	assert.Equal(t, "someMessage", logs[0].Message)
	assert.ElementsMatch(t, []zap.Field{
		zap.String("someString", "someStringValue"),
		zap.Int("someInt", 52),
		zap.Duration("someDuration", timeSince),
		zap.Any("someAny", map[string]interface{}{"someKey": "someValue"}),
	}, logs[0].Context)
}
