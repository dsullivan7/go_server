package routes

import (
  "github.com/go-chi/chi"
  "github.com/go-chi/chi/middleware"

  UserRoutes "go_server/internal/routes/users"
)

func Init() *chi.Mux{
	router := chi.NewRouter()

  router.Use(middleware.RequestID)
  router.Use(middleware.RealIP)
  router.Use(middleware.Logger)
  router.Use(middleware.Recoverer)

  router.Mount("/api", initApi())

	return router
}

func initApi() *chi.Mux{
  router := chi.NewRouter()

  router.Mount("/users", UserRoutes.Routes())

	return router
}
