package server

import (
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (server *Server) initRoutes() *chi.Mux {
	router := server.router
	controllers := server.controllers

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(server.middlewares.Logger())
	router.Use(server.middlewares.ContentType("application/json; charset=utf-8"))
	router.Use(server.middlewares.HandlePanic())

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(server.config.AllowedOrigins, ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           server.config.RouterMaxAge,
	}))

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

	return router
}
