package controllers

import (
	"go_server/internal/config"
	"go_server/internal/logger"
	"go_server/internal/store"
)

const HTTP400 = 400

type Controllers struct {
	store  store.Store
	config *config.Config
	logger logger.Logger
}

func NewControllers(store store.Store, config *config.Config, logger logger.Logger) *Controllers {
	return &Controllers{
		store:  store,
		config: config,
		logger: logger,
	}
}
