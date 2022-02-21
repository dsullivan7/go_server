package controllers

import (
	"net/http"
	"go_server/internal/errors"
	"github.com/google/uuid"

	"github.com/go-chi/render"
)

func (c *Controllers) ListPortfolioHoldings(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	portfolioID := uuid.Must(uuid.Parse(r.URL.Query().Get("portfolio_id")))

	query["portfolio_id"] = portfolioID

	portfolio, errPortfolio := c.store.GetPortfolio(portfolioID)

	if errPortfolio != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errPortfolio})

		return
	}

	securities, errSecurities := c.store.ListSecurities(map[string]interface{}{})

	if errSecurities != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errSecurities})

		return
	}

	securityTags, errSecurityTags := c.store.ListSecurityTags(map[string]interface{}{})

	if errSecurityTags != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errSecurityTags})

		return
	}

	portfolioTags, errPortfolioTags := c.store.ListPortfolioTags(query)

	if errPortfolioTags != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errPortfolioTags})

		return
	}

	portfolioHoldings := c.services.GetPortfolioHoldings(*portfolio, portfolioTags, securities, securityTags)

	render.JSON(w, r, portfolioHoldings)
}
