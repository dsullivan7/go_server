package logger

import (
	"time"

	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(logger *zap.Logger) Logger {
	return &ZapLogger{logger: logger}
}

func (logger *ZapLogger) Info(message string) {
	logger.logger.Info(message)
}

func (logger *ZapLogger) InfoWithMeta(message string, meta map[string]interface{}) {
	var args []zap.Field

	for key, value := range meta {
		if valueWithType, ok := value.(string); ok {
			args = append(args, zap.String(key, valueWithType))
		} else if valueWithType, ok := value.(int); ok {
			args = append(args, zap.Int(key, valueWithType))
		} else if valueWithType, ok := value.(time.Duration); ok {
			args = append(args, zap.Duration(key, valueWithType))
		} else {
			args = append(args, zap.Any(key, value))
		}
	}

	logger.logger.Info(message, args...)
}
