package users

import (
  "github.com/go-chi/chi"

  UserController "go_server/internal/controllers/users"
)

func Routes() *chi.Mux {
  router := chi.NewRouter()
  router.Get("/", UserController.Get)
  return router
}
