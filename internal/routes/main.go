package routes

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	customMiddlewares "go_server/internal/middlewares"
	"go_server/internal/logger"
	"go_server/internal/config"
)

func Init() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(customMiddlewares.Logger(logger.Logger))
	router.Use(middleware.Recoverer)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(config.AllowedOrigins, ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Mount("/api", initAPI())
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	return router
}

func initAPI() *chi.Mux {
	router := chi.NewRouter()

	// router.With(GoServerMiddlewares.Auth).With(GoServerMiddlewares.User).Mount("/users", UserRoutes())
	router.Mount("/users", UserRoutes())
	router.Mount("/reviews", ReviewRoutes())

	return router
}
