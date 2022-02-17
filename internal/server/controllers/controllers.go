package controllers

import (
	"go_server/internal/broker"
	"go_server/internal/config"
	"go_server/internal/logger"
	"go_server/internal/plaid"
	"go_server/internal/server/utils"
	"go_server/internal/store"
)

type Controllers struct {
	config      *config.Config
	store       store.Store
	plaidClient plaid.IClient
	broker      broker.Broker
	utils       *utils.ServerUtils
	logger      logger.Logger
}

func NewControllers(
	cfg *config.Config,
	str store.Store,
	plaidClient plaid.IClient,
	brkr broker.Broker,
	utls *utils.ServerUtils,
	lggr logger.Logger,
) *Controllers {
	return &Controllers{
		store:       str,
		config:      cfg,
		plaidClient: plaidClient,
		broker:      brkr,
		logger:      lggr,
		utils:       utls,
	}
}
