package logger

import (
  "go.uber.org/zap"
)

type ZapLogger struct {
  logger *zap.Logger
}

func NewZapLogger() Logger {
  logger, _ := zap.NewProduction()

  return &ZapLogger{ logger: logger }
}

func (l *ZapLogger) Info(args ...interface{}) {
	println("blah")
}
