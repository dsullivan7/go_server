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

func (c *Controllers) GetPortfolioTag(w http.ResponseWriter, r *http.Request) {
	portfolioTagID := uuid.Must(uuid.Parse(chi.URLParam(r, "portfolio_tag_id")))

	portfolioTag, err := c.store.GetPortfolioTag(portfolioTagID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, portfolioTag)
}

func (c *Controllers) ListPortfolioTags(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	portfolioID := r.URL.Query().Get("portfolio_id")

	if portfolioID != "" {
		query["portfolio_id"] = portfolioID
	}

	portfolioTags, err := c.store.ListPortfolioTags(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, portfolioTags)
}

func (c *Controllers) CreatePortfolioTag(w http.ResponseWriter, r *http.Request) {
	var portfolioTagPayload models.PortfolioTag

	errDecode := json.NewDecoder(r.Body).Decode(&portfolioTagPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	portfolioTag, err := c.store.CreatePortfolioTag(portfolioTagPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, portfolioTag)
}

func (c *Controllers) ModifyPortfolioTag(w http.ResponseWriter, r *http.Request) {
	var portfolioTagPayload models.PortfolioTag

	portfolioTagID := uuid.Must(uuid.Parse(chi.URLParam(r, "portfolio_tag_id")))

	errDecode := json.NewDecoder(r.Body).Decode(&portfolioTagPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	portfolioTag, err := c.store.ModifyPortfolioTag(portfolioTagID, portfolioTagPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, portfolioTag)
}

func (c *Controllers) DeletePortfolioTag(w http.ResponseWriter, r *http.Request) {
	portfolioTagID := uuid.Must(uuid.Parse(chi.URLParam(r, "portfolio_tag_id")))

	err := c.store.DeletePortfolioTag(portfolioTagID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.NoContent(w, r)
}
