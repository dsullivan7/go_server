package controllers

import (
	"go_server/internal/bank"
	"go_server/internal/config"
	"go_server/internal/crawler"
	"go_server/internal/logger"
	"go_server/internal/server/utils"
	"go_server/internal/store"
)

type Controllers struct {
	config  *config.Config
	store   store.Store
	crawler crawler.Crawler
	bank    bank.Bank
	utils   *utils.ServerUtils
	logger  logger.Logger
}

func NewControllers(
	cfg *config.Config,
	str store.Store,
	crwlr crawler.Crawler,
	bnk bank.Bank,
	utls *utils.ServerUtils,
	lggr logger.Logger,
) *Controllers {
	return &Controllers{
		store:   str,
		config:  cfg,
		crawler: crwlr,
		bank:    bnk,
		logger:  lggr,
		utils:   utls,
	}
}
