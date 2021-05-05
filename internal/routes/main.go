package routes

import (
	UserRoutes "go_server/internal/routes/users"

	GoServerMiddlewares "go_server/internal/middlewares"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Init() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Mount("/api", initAPI())

	return router
}

func initAPI() *chi.Mux {
	router := chi.NewRouter()

	router.With(GoServerMiddlewares.Auth).Mount("/users", UserRoutes.Routes())

	return router
}
