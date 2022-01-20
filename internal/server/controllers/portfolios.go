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

func (c *Controllers) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	portfolioID := uuid.Must(uuid.Parse(chi.URLParam(r, "portfolio_id")))

	portfolio, err := c.store.GetPortfolio(portfolioID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, portfolio)
}

func (c *Controllers) ListPortfolios(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	fromUserID := r.URL.Query().Get("from_user_id")
	toUserID := r.URL.Query().Get("to_user_id")

	if fromUserID != "" {
		query["from_user_id"] = fromUserID
	}

	if toUserID != "" {
		query["to_user_id"] = toUserID
	}

	portfolios, err := c.store.ListPortfolios(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, portfolios)
}

func (c *Controllers) CreatePortfolio(w http.ResponseWriter, r *http.Request) {
	var portfolioPayload models.Portfolio

	errDecode := json.NewDecoder(r.Body).Decode(&portfolioPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	portfolio, err := c.store.CreatePortfolio(portfolioPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, portfolio)
}

func (c *Controllers) ModifyPortfolio(w http.ResponseWriter, r *http.Request) {
	var portfolioPayload models.Portfolio

	portfolioID := uuid.Must(uuid.Parse(chi.URLParam(r, "portfolio_id")))

	errDecode := json.NewDecoder(r.Body).Decode(&portfolioPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	portfolio, err := c.store.ModifyPortfolio(portfolioID, portfolioPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, portfolio)
}

func (c *Controllers) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	portfolioID := uuid.Must(uuid.Parse(chi.URLParam(r, "portfolio_id")))

	err := c.store.DeletePortfolio(portfolioID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
