package server

import (
	"go_server/internal/server/graph/generated"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func (s *ChiServer) Init() http.Handler {
	router := s.router
	controllers := s.controllers
	resolver := s.resolver
	middlewares := s.middlewares

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(render.SetContentType(render.ContentTypeJSON))
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

		r.Get("/api/users/{user_id}", controllers.GetUser)
		r.Get("/api/users", controllers.ListUsers)
		r.Post("/api/users", controllers.CreateUser)
		r.Delete("/api/users/{user_id}", controllers.DeleteUser)
		r.Put("/api/users/{user_id}", controllers.ModifyUser)

		r.Get("/api/reviews/{review_id}", controllers.GetReview)
		r.Get("/api/reviews", controllers.ListReviews)
		r.Post("/api/reviews", controllers.CreateReview)
		r.Delete("/api/reviews/{review_id}", controllers.DeleteReview)
		r.Put("/api/reviews/{review_id}", controllers.ModifyReview)

		r.Get("/api/industries/{industry_id}", controllers.GetIndustry)
		r.Get("/api/industries", controllers.ListIndustries)
		r.Post("/api/industries", controllers.CreateIndustry)
		r.Delete("/api/industries/{industry_id}", controllers.DeleteIndustry)
		r.Put("/api/industries/{industry_id}", controllers.ModifyIndustry)

		r.Get("/api/bank-accounts/{bank_account_id}", controllers.GetBankAccount)
		r.Get("/api/bank-accounts", controllers.ListBankAccounts)
		r.Post("/api/bank-accounts", controllers.CreateBankAccount)
		r.Delete("/api/bank-accounts/{bank_account_id}", controllers.DeleteBankAccount)
		r.Put("/api/bank-accounts/{bank_account_id}", controllers.ModifyBankAccount)

		// r.Get("/api/brokerage-accounts/{brokerage_account_id}", controllers.GetBrokerageAccount)
		// r.Get("/api/brokerage-accounts", controllers.ListBrokerageAccounts)
		// r.Post("/api/brokerage-accounts", controllers.CreateBrokerageAccount)
		// r.Delete("/api/brokerage-accounts/{brokerage_account_id}", controllers.DeleteBrokerageAccount)
		// r.Put("/api/brokerage-accounts/{brokerage_account_id}", controllers.ModifyBrokerageAccount)

		r.Get("/api/portfolios/{portfolio_id}", controllers.GetPortfolio)
		r.Get("/api/portfolios", controllers.ListPortfolios)
		r.Post("/api/portfolios", controllers.CreatePortfolio)
		r.Delete("/api/portfolios/{portfolio_id}", controllers.DeletePortfolio)
		r.Put("/api/portfolios/{portfolio_id}", controllers.ModifyPortfolio)

		r.Get("/api/portfolio-industries/{portfolio_industry_id}", controllers.GetPortfolioIndustry)
		r.Get("/api/portfolio-industries", controllers.ListPortfolioIndustries)
		r.Post("/api/portfolio-industries", controllers.CreatePortfolioIndustry)
		r.Delete("/api/portfolio-industries/{portfolio_industry_id}", controllers.DeletePortfolioIndustry)
		r.Put("/api/portfolio-industries/{portfolio_industry_id}", controllers.ModifyPortfolioIndustry)

		r.Post("/api/plaid/token", controllers.CreatePlaidToken)

		r.Get("/api/snap", controllers.GetSnap)

		handler := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
		r.Handle("/query", handler)
	})

	return router
}
