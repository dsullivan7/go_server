package users

import (
  "github.com/go-chi/chi"

  UsersController "go_server/internal/controllers/users"
)

func Routes() *chi.Mux {
  router := chi.NewRouter()
  router.Get("/{userId}", UsersController.Get)
  return router
}
