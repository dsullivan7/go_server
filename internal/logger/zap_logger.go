package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
}

type zapLog func(string, ...zapcore.Field)

func NewZapLogger(logger *zap.Logger) Logger {
	return &ZapLogger{logger: logger}
}

func (logger *ZapLogger) log(fn zapLog, message string) {
	fn(message)
}

func (logger *ZapLogger) logWithMeta(fn zapLog, message string, meta map[string]interface{}) {
	var args []zap.Field

	for key, value := range meta {
		if valueWithType, ok := value.(string); ok {
			args = append(args, zap.String(key, valueWithType))
		} else if valueWithType, ok := value.(int); ok {
			args = append(args, zap.Int(key, valueWithType))
		} else if valueWithType, ok := value.(time.Duration); ok {
			args = append(args, zap.Duration(key, valueWithType))
		} else if valueWithType, ok := value.(error); ok {
			args = append(args, zap.Error(valueWithType))
		} else {
			args = append(args, zap.Any(key, value))
		}
	}

	fn(message, args...)
}

func (logger *ZapLogger) Info(message string) {
	logger.log(logger.logger.Info, message)
}

func (logger *ZapLogger) InfoWithMeta(message string, meta map[string]interface{}) {
	logger.logWithMeta(logger.logger.Info, message, meta)
}

func (logger *ZapLogger) Error(message string) {
	logger.log(logger.logger.Error, message)
}

func (logger *ZapLogger) ErrorWithMeta(message string, meta map[string]interface{}) {
	logger.logWithMeta(logger.logger.Error, message, meta)
}
