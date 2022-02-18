package controllers

import (
	"go_server/internal/errors"
	"go_server/internal/models"
	"net/http"
	"github.com/google/uuid"

	"github.com/go-chi/render"
)

var securityID1 = uuid.New()

var securities = []models.Security{
	models.Security{
		Symbol: "APPL",
		SecurityID: securityID1,
	},
}

var securityTags = []models.SecurityTag{
	{
		SecurityID: securityID1,
		SecurityID: securityID1,
	},
}

func (c *Controllers) ListPortfolioHoldings(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	portfolioID := r.URL.Query().Get("portfolio_id")

	if portfolioID != "" {
		query["portfolio_id"] = portfolioID
	}

	securities, errSecurities := c.store.ListSecurities(query)

	if errSecurities != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errSecurities})

		return
	}

	securityTags, errSecurityTags := c.store.ListSecurityTags(query)

	if errSecurities != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errSecurities})

		return
	}

	portfolioTags, errPortfolioTags := c.store.ListPortfolioTags(query)

	if errPortfolioTags != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errPortfolioTags})

		return
	}

	c.services.GetPortfolio()

	render.JSON(w, r, portfolioHoldings)
}
