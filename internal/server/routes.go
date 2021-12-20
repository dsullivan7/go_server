package server

import (
	"go_server/internal/server/graph/generated"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (s *ChiServer) Init() http.Handler {
	router := s.router
	controllers := s.controllers
	resolver := s.resolver
	middlewares := s.middlewares

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.ContentType("application/json; charset=utf-8"))
	router.Use(middlewares.Logger())
	router.Use(middlewares.HandlePanic())

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   s.config.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           s.config.RouterMaxAge,
	}))

	router.Get("/", controllers.GetHealth)

	router.Group(func(r chi.Router) {
		r.Use(middlewares.Auth())
		r.Use(middlewares.User())

		r.Get("/api/users/{userID}", controllers.GetUser)
		r.Get("/api/users", controllers.ListUsers)
		r.Post("/api/users", controllers.CreateUser)
		r.Delete("/api/users/{userID}", controllers.DeleteUser)
		r.Put("/api/users/{userID}", controllers.ModifyUser)

		r.Get("/api/reviews/{reviewID}", controllers.GetReview)
		r.Get("/api/reviews", controllers.ListReviews)
		r.Post("/api/reviews", controllers.CreateReview)
		r.Delete("/api/reviews/{reviewID}", controllers.DeleteReview)
		r.Put("/api/reviews/{reviewID}", controllers.ModifyReview)

		r.Get("/api/snap", controllers.GetSnap)

		handler := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
		r.Handle("/query", handler)
	})

	return router
}
