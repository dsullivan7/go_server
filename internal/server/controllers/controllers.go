package controllers

import (
	"go_server/internal/config"
	"go_server/internal/logger"
	"go_server/internal/server/utils"
	"go_server/internal/store"
)

type Controllers struct {
	store  store.Store
	config *config.Config
	logger logger.Logger
	utils  *utils.ServerUtils
}

func NewControllers(
	store store.Store,
	config *config.Config,
	logger logger.Logger,
	utils *utils.ServerUtils,
) *Controllers {
	return &Controllers{
		store:  store,
		config: config,
		logger: logger,
		utils:  utils,
	}
}
