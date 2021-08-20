package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() (Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize zap logger: %w", err)
	}

	return &ZapLogger{logger: logger}, nil
}

func (logger *ZapLogger) Info(message string) {
	logger.logger.Info(message)
}
