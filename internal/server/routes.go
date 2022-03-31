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

	router.Post("/api/credentials", controllers.CreateCredential)

	router.Post("/api/plaid/token", controllers.CreatePlaidToken)

	router.Get("/api/invoices/{invoice_id}", controllers.GetInvoice)
	router.Get("/api/invoices", controllers.ListInvoices)
	router.Post("/api/invoices", controllers.CreateInvoice)
	router.Delete("/api/invoices/{invoice_id}", controllers.DeleteInvoice)
	router.Put("/api/invoices/{invoice_id}", controllers.ModifyInvoice)

	router.Get("/api/items/{item_id}", controllers.GetItem)
	router.Get("/api/items", controllers.ListItems)
	router.Post("/api/items", controllers.CreateItem)
	router.Delete("/api/items/{item_id}", controllers.DeleteItem)
	router.Put("/api/items/{item_id}", controllers.ModifyItem)

	router.Post("/api/bank-accounts", controllers.CreateBankAccount)

	router.Group(func(r chi.Router) {
		r.Use(middlewares.Auth())
		r.Use(middlewares.User())

		r.Get("/api/users/{user_id}", controllers.GetUser)
		r.Get("/api/users", controllers.ListUsers)
		r.Post("/api/users", controllers.CreateUser)
		r.Delete("/api/users/{user_id}", controllers.DeleteUser)
		r.Put("/api/users/{user_id}", controllers.ModifyUser)

		r.Get("/api/tags/{tag_id}", controllers.GetTag)
		r.Get("/api/tags", controllers.ListTags)
		r.Post("/api/tags", controllers.CreateTag)
		r.Delete("/api/tags/{tag_id}", controllers.DeleteTag)
		r.Put("/api/tags/{tag_id}", controllers.ModifyTag)

		r.Get("/api/bank-accounts/{bank_account_id}", controllers.GetBankAccount)
		r.Get("/api/bank-accounts", controllers.ListBankAccounts)
		r.Delete("/api/bank-accounts/{bank_account_id}", controllers.DeleteBankAccount)
		r.Put("/api/bank-accounts/{bank_account_id}", controllers.ModifyBankAccount)

		r.Get("/api/bank-transfers/{bank_transfer_id}", controllers.GetBankTransfer)
		r.Get("/api/bank-transfers", controllers.ListBankTransfers)
		r.Post("/api/bank-transfers", controllers.CreateBankTransfer)
		r.Delete("/api/bank-transfers/{bank_transfer_id}", controllers.DeleteBankTransfer)
		r.Put("/api/bank-transfers/{bank_transfer_id}", controllers.ModifyBankTransfer)

		r.Get("/api/brokerage-accounts/{brokerage_account_id}", controllers.GetBrokerageAccount)
		r.Get("/api/brokerage-accounts", controllers.ListBrokerageAccounts)
		r.Post("/api/brokerage-accounts", controllers.CreateBrokerageAccount)
		r.Delete("/api/brokerage-accounts/{brokerage_account_id}", controllers.DeleteBrokerageAccount)
		r.Put("/api/brokerage-accounts/{brokerage_account_id}", controllers.ModifyBrokerageAccount)

		r.Get("/api/portfolios/{portfolio_id}", controllers.GetPortfolio)
		r.Get("/api/portfolios", controllers.ListPortfolios)
		r.Post("/api/portfolios", controllers.CreatePortfolio)
		r.Delete("/api/portfolios/{portfolio_id}", controllers.DeletePortfolio)
		r.Put("/api/portfolios/{portfolio_id}", controllers.ModifyPortfolio)

		r.Get("/api/portfolio-tags/{portfolio_tag_id}", controllers.GetPortfolioTag)
		r.Get("/api/portfolio-tags", controllers.ListPortfolioTags)
		r.Post("/api/portfolio-tags", controllers.CreatePortfolioTag)
		r.Delete("/api/portfolio-tags/{portfolio_tag_id}", controllers.DeletePortfolioTag)
		r.Put("/api/portfolio-tags/{portfolio_tag_id}", controllers.ModifyPortfolioTag)

		r.Get("/api/tags/{tag_id}", controllers.GetTag)
		r.Get("/api/tags", controllers.ListTags)
		r.Post("/api/tags", controllers.CreateTag)
		r.Delete("/api/tags/{tag_id}", controllers.DeleteTag)
		r.Put("/api/tags/{tag_id}", controllers.ModifyTag)

		r.Get("/api/orders/{order_id}", controllers.GetOrder)
		r.Get("/api/orders", controllers.ListOrders)
		r.Post("/api/orders", controllers.CreateOrder)
		r.Delete("/api/orders/{order_id}", controllers.DeleteOrder)
		r.Put("/api/orders/{order_id}", controllers.ModifyOrder)

		r.Get("/api/portfolio-recommendations", controllers.ListPortfolioRecommendations)

		r.Get("/api/positions", controllers.ListPositions)

		handler := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
		r.Handle("/query", handler)
	})

	return router
}
