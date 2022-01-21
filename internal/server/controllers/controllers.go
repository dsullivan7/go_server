package controllers

import (
	"go_server/internal/config"
	"go_server/internal/crawler"
	"go_server/internal/logger"
	"go_server/internal/plaid"
	"go_server/internal/server/utils"
	"go_server/internal/store"
)

type Controllers struct {
	config      *config.Config
	store       store.Store
	crawler     crawler.Crawler
	plaidClient plaid.Client
	utils       *utils.ServerUtils
	logger      logger.Logger
}

func NewControllers(
	config *config.Config,
	store store.Store,
	crawler crawler.Crawler,
	plaidClient plaid.Client,
	utils *utils.ServerUtils,
	logger logger.Logger,
) *Controllers {
	return &Controllers{
		store:       store,
		config:      config,
		crawler:     crawler,
		plaidClient: plaidClient,
		logger:      logger,
		utils:       utils,
	}
}
