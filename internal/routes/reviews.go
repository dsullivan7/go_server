package routes

import (
	"go_server/internal/controllers"

	"github.com/go-chi/chi"
)

func ReviewRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{reviewID}", controllers.GetReview)
	router.Get("/", controllers.ListReviews)
	router.Post("/", controllers.CreateReview)
	router.Delete("/{reviewID}", controllers.DeleteReview)
	router.Put("/{reviewID}", controllers.ModifyReview)

	return router
}
