package routes

import (
  "github.com/go-chi/chi"

  UserRoutes "go_server/internal/routes/users"
)

func Init() *chi.Mux{
	router := chi.NewRouter()

  router.Mount("/api/users", UserRoutes.Routes())

	return router
}
