package controllers

import (
	"go_server/internal/config"
	"go_server/internal/store"
	"go_server/internal/logger"
)

type Controllers struct {
	store store.Store
  config *config.Config
	logger logger.Logger
}

func NewControllers(store store.Store, config *config.Config, logger logger.Logger) *Controllers {
  return &Controllers{
		store: store,
    config: config,
		logger: logger,
	}
}
