package users

import (
	UsersController "go_server/internal/controllers/users"

	"github.com/go-chi/chi"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{userID}", UsersController.Get)
	router.Get("/", UsersController.List)
	router.Post("/", UsersController.Create)
	router.Delete("/{userID}", UsersController.Delete)
	router.Put("/{userID}", UsersController.Modify)

	return router
}
