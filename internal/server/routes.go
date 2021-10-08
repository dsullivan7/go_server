package server

import (
	"net/http"
	"github.com/99designs/gqlgen/graphql/handler"
	"go_server/internal/server/graph/generated"

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

	router.Use(middlewares.Auth())
	router.Use(middlewares.User())

	router.Get("/api/users/{userID}", controllers.GetUser)
	router.Get("/api/users", controllers.ListUsers)
	router.Post("/api/users", controllers.CreateUser)
	router.Delete("/api/users/{userID}", controllers.DeleteUser)
	router.Put("/api/users/{userID}", controllers.ModifyUser)

	router.Get("/api/reviews/{reviewID}", controllers.GetReview)
	router.Get("/api/reviews", controllers.ListReviews)
	router.Post("/api/reviews", controllers.CreateReview)
	router.Delete("/api/reviews/{reviewID}", controllers.DeleteReview)
	router.Put("/api/reviews/{reviewID}", controllers.ModifyReview)

	handler := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver }))
	router.Handle("/query", handler)

	return router
}
