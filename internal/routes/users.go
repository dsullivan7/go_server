package routes

import (
	"go_server/internal/controllers"

	"github.com/go-chi/chi"
)

func UserRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{userID}", controllers.GetUser)
	router.Get("/", controllers.ListUsers)
	router.Post("/", controllers.CreateUser)
	router.Delete("/{userID}", controllers.DeleteUser)
	router.Put("/{userID}", controllers.ModifyUser)

	return router
}
