package middlewares

import (
	"go_server/internal/authentication"
	"go_server/internal/config"
	"go_server/internal/logger"
	"go_server/internal/server/utils"
	"go_server/internal/store"
)

type Middlewares struct {
	config *config.Config
	store  store.Store
	auth   authentication.Authentication
	logger logger.Logger
	utils  *utils.ServerUtils
}

func NewMiddlewares(
	config *config.Config,
	store store.Store,
	ath authentication.Authentication,
	utils *utils.ServerUtils,
	logger logger.Logger,
) *Middlewares {
	return &Middlewares{
		config: config,
		store:  store,
		auth:   ath,
		utils:  utils,
		logger: logger,
	}
}
