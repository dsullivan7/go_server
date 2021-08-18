package server

import (
	"go_server/internal/controllers"
	"go_server/internal/logger"
	"go_server/internal/config"
	"github.com/go-chi/chi"
	"fmt"
	"net/http"
)

type Server struct {
	router *chi.Mux
	controllers *controllers.Controllers
	logger logger.Logger
	config *config.Config
}

func NewServer(
	router *chi.Mux,
	controllers *controllers.Controllers,
	config *config.Config,
	logger logger.Logger,
) *Server {
	return &Server{
		router: router,
		controllers: controllers,
		config: config,
		logger: logger,
	}
}

func (server *Server) Run() {
	server.Routes()

  server.logger.Info(fmt.Sprintf("Server started on port %s", server.config.Port))

  server.logger.Info(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), server.router))
}
