package controllers

import (
	"go_server/internal/bank"
	"go_server/internal/broker"
	"go_server/internal/cipher"
	"go_server/internal/config"
	"go_server/internal/gov"
	"go_server/internal/logger"
	"go_server/internal/plaid"
	"go_server/internal/server/utils"
	"go_server/internal/services"
	"go_server/internal/store"
)

type Controllers struct {
	config      *config.Config
	store       store.Store
	plaidClient plaid.IClient
	services    services.IService
	cipher      cipher.ICipher
	bank        bank.Bank
	broker      broker.Broker
	gov         gov.IGov
	utils       *utils.ServerUtils
	logger      logger.Logger
}

func NewControllers(
	cfg *config.Config,
	srvc services.IService,
	str store.Store,
	plaidClient plaid.IClient,
	brkr broker.Broker,
	bnk bank.Bank,
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
		bank:        bnk,
		gov:         gv,
		logger:      lggr,
		utils:       utls,
	}
}
