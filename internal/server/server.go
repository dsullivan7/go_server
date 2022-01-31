package server

import (
	"net/http"

	"go_server/internal/auth"
	"go_server/internal/config"
	"go_server/internal/crawler"
	"go_server/internal/logger"
	"go_server/internal/bank"
	"go_server/internal/server/controllers"
	"go_server/internal/server/graph"
	"go_server/internal/server/middlewares"
	"go_server/internal/server/utils"
	"go_server/internal/store"

	"github.com/go-chi/chi"
)

type Server interface {
	Init() http.Handler
}

type ChiServer struct {
	router      *chi.Mux
	config      *config.Config
	controllers *controllers.Controllers
	resolver    *graph.Resolver
	middlewares *middlewares.Middlewares
	logger      logger.Logger
}

func NewChiServer(
	cfg *config.Config,
	router *chi.Mux,
	str store.Store,
	crwlr crawler.Crawler,
	bnk bank.Bank,
	ath auth.Auth,
	lggr logger.Logger,
) Server {
	utils := utils.NewServerUtils(lggr)
	controllers := controllers.NewControllers(cfg, str, crwlr, bnk, utils, lggr)
	resolver := graph.NewResolver(cfg, str, lggr)
	middlewares := middlewares.NewMiddlewares(cfg, str, ath, utils, lggr)

	return &ChiServer{
		router:      router,
		config:      cfg,
		logger:      lggr,
		controllers: controllers,
		resolver:    resolver,
		middlewares: middlewares,
	}
}
