package zap

import (
	"time"

	"go_server/internal/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

type zapLog func(string, ...zapcore.Field)

func NewLogger(logger *zap.Logger) logger.Logger {
	return &Logger{logger: logger}
}

func (logger *Logger) log(fn zapLog, message string) {
	fn(message)
}

func (logger *Logger) logWithMeta(fn zapLog, message string, meta map[string]interface{}) {
	var args []zap.Field

	for key, value := range meta {
		switch valueWithType := value.(type) {
		case string:
			args = append(args, zap.String(key, valueWithType))
		case int:
			args = append(args, zap.Int(key, valueWithType))
		case time.Duration:
			args = append(args, zap.Duration(key, valueWithType))
		case error:
			args = append(args, zap.Error(valueWithType))
		default:
			args = append(args, zap.Any(key, value))
		}
	}

	fn(message, args...)
}

func (logger *Logger) Info(message string) {
	logger.log(logger.logger.Info, message)
}

func (logger *Logger) InfoWithMeta(message string, meta map[string]interface{}) {
	logger.logWithMeta(logger.logger.Info, message, meta)
}

func (logger *Logger) Error(message string) {
	logger.log(logger.logger.Error, message)
}

func (logger *Logger) ErrorWithMeta(message string, meta map[string]interface{}) {
	logger.logWithMeta(logger.logger.Error, message, meta)
}
