package controllers

import (
	"encoding/json"
	"go_server/internal/errors"
	"go_server/internal/models"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func (c *Controllers) GetPortfolioIndustry(w http.ResponseWriter, r *http.Request) {
	portfolioIndustryID := uuid.Must(uuid.Parse(chi.URLParam(r, "portfolio_industry_id")))

	portfolioIndustry, err := c.store.GetPortfolioIndustry(portfolioIndustryID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, portfolioIndustry)
}

func (c *Controllers) ListPortfolioIndustries(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	portfolioID := r.URL.Query().Get("portfolio_id")

	if portfolioID != "" {
		query["portfolio_id"] = portfolioID
	}

	portfolioIndustries, err := c.store.ListPortfolioIndustries(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, portfolioIndustries)
}

func (c *Controllers) CreatePortfolioIndustry(w http.ResponseWriter, r *http.Request) {
	var portfolioIndustryPayload models.PortfolioIndustry

	errDecode := json.NewDecoder(r.Body).Decode(&portfolioIndustryPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	portfolioIndustry, err := c.store.CreatePortfolioIndustry(portfolioIndustryPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, portfolioIndustry)
}

func (c *Controllers) ModifyPortfolioIndustry(w http.ResponseWriter, r *http.Request) {
	var portfolioIndustryPayload models.PortfolioIndustry

	portfolioIndustryID := uuid.Must(uuid.Parse(chi.URLParam(r, "portfolio_industry_id")))

	errDecode := json.NewDecoder(r.Body).Decode(&portfolioIndustryPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	portfolioIndustry, err := c.store.ModifyPortfolioIndustry(portfolioIndustryID, portfolioIndustryPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, portfolioIndustry)
}

func (c *Controllers) DeletePortfolioIndustry(w http.ResponseWriter, r *http.Request) {
	portfolioIndustryID := uuid.Must(uuid.Parse(chi.URLParam(r, "portfolio_industry_id")))

	err := c.store.DeletePortfolioIndustry(portfolioIndustryID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
