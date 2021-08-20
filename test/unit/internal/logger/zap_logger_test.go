package logger_test

import (
	"go_server/internal/logger"
	"testing"
	"encoding/json"

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

	logger.InfoWithMeta(
		"blah",
		map[string]interface{}{
			"some": "key",
		},
	)

	logs := recordedLogs.All()

	assert.Equal(t, 1, len(logs))

	var logJSON map[string]string

	errDecode := json.Unmarshal([]byte(`{"some":"key"}`), &logJSON)
	assert.Nil(t, errDecode)

	assert.Equal(t, "key", logJSON["some"])
}
