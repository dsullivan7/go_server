package middlewares

import (
	"go_server/internal/config"
	"go_server/internal/logger"
	"go_server/internal/server/utils"
	"go_server/internal/store"
)

type Middlewares struct {
	store  store.Store
	config *config.Config
	logger logger.Logger
	utils  *utils.ServerUtils
}

func NewMiddlewares(
	store store.Store,
	config *config.Config,
	logger logger.Logger,
	utils *utils.ServerUtils,
) *Middlewares {
	return &Middlewares{
		store:  store,
		config: config,
		logger: logger,
		utils:  utils,
	}
}
