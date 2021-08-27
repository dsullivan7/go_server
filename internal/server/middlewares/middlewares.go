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
	config *config.Config,
	store store.Store,
	utils *utils.ServerUtils,
	logger logger.Logger,
) *Middlewares {
	return &Middlewares{
		config: config,
		store:  store,
		utils:  utils,
		logger: logger,
	}
}
