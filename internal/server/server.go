package server

import (
	"fmt"
	"go_server/internal/config"
	"go_server/internal/logger"
	"go_server/internal/server/controllers"
	"go_server/internal/server/middlewares"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
	router      *chi.Mux
	controllers *controllers.Controllers
	middlewares *middlewares.Middlewares
	logger      logger.Logger
	config      *config.Config
}

func NewServer(
	router *chi.Mux,
	controllers *controllers.Controllers,
	middlewares *middlewares.Middlewares,
	config *config.Config,
	logger logger.Logger,
) *Server {
	return &Server{
		router:      router,
		controllers: controllers,
		middlewares: middlewares,
		config:      config,
		logger:      logger,
	}
}

func (server *Server) Run() {
	server.Routes()

	server.logger.Info(fmt.Sprintf("Server started on port %s", server.config.Port))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), server.router))
}
