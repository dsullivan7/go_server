package controllers

import (
	"go_server/internal/broker"
	"go_server/internal/cipher"
	"go_server/internal/config"
	"go_server/internal/logger"
	"go_server/internal/plaid"
	"go_server/internal/server/utils"
	"go_server/internal/services"
	"go_server/internal/store"
	"go_server/internal/gov"
)

type Controllers struct {
	config      *config.Config
	store       store.Store
	plaidClient plaid.IClient
	services    services.IService
	cipher      cipher.ICipher
	broker      broker.Broker
	gov      gov.IGov
	utils       *utils.ServerUtils
	logger      logger.Logger
}

func NewControllers(
	cfg *config.Config,
	srvc services.IService,
	str store.Store,
	plaidClient plaid.IClient,
	brkr broker.Broker,
	gv gov.IGov,
	cphr cipher.ICipher,
	utls *utils.ServerUtils,
	lggr logger.Logger,
) *Controllers {
	return &Controllers{
		store:       str,
		services:    srvc,
		config:      cfg,
		cipher:      cphr,
		plaidClient: plaidClient,
		broker:      brkr,
		gov:      gv,
		logger:      lggr,
		utils:       utls,
	}
}
