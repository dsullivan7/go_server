package zap_test

import (
	"errors"
	goServerZapLogger "go_server/internal/logger/zap"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var errTest = errors.New("this is an error")

func TestZapLoggerInfo(t *testing.T) {
	t.Parallel()

	core, recordedLogs := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(core)

	logger := goServerZapLogger.NewLogger(zapLogger)

	logger.Info("someMessage")

	logs := recordedLogs.All()

	assert.Equal(t, 1, len(logs))
	assert.Equal(t, zapcore.Level(0), logs[0].Level)
	assert.Equal(t, "someMessage", logs[0].Message)
}

func TestZapLoggerInfoWithMeta(t *testing.T) {
	t.Parallel()

	core, recordedLogs := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(core)

	logger := goServerZapLogger.NewLogger(zapLogger)

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
	assert.Equal(t, zapcore.Level(0), logs[0].Level)
	assert.Equal(t, "someMessage", logs[0].Message)
	assert.ElementsMatch(t, []zap.Field{
		zap.String("someString", "someStringValue"),
		zap.Int("someInt", 52),
		zap.Duration("someDuration", timeSince),
		zap.Any("someAny", map[string]interface{}{"someKey": "someValue"}),
	}, logs[0].Context)
}

func TestZapLoggerError(t *testing.T) {
	t.Parallel()

	core, recordedLogs := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(core)

	logger := goServerZapLogger.NewLogger(zapLogger)

	logger.Error("someError")

	logs := recordedLogs.All()

	assert.Equal(t, 1, len(logs))
	assert.Equal(t, zapcore.Level(2), logs[0].Level)
	assert.Equal(t, "someError", logs[0].Message)
}

func TestZapLoggerErrorWithMeta(t *testing.T) {
	t.Parallel()

	core, recordedLogs := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(core)

	logger := goServerZapLogger.NewLogger(zapLogger)

	logger.ErrorWithMeta(
		"someError",
		map[string]interface{}{
			"someError": errTest,
		},
	)

	logs := recordedLogs.All()

	assert.Equal(t, 1, len(logs))
	assert.Equal(t, zapcore.Level(2), logs[0].Level)
	assert.Equal(t, "someError", logs[0].Message)
	assert.ElementsMatch(t, []zap.Field{
		zap.Error(errTest),
	}, logs[0].Context)
}
