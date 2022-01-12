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

func (c *Controllers) GetIndustry(w http.ResponseWriter, r *http.Request) {
	industryID := uuid.Must(uuid.Parse(chi.URLParam(r, "industryID")))

	industry, err := c.store.GetIndustry(industryID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, industry)
}

func (c *Controllers) ListIndustries(w http.ResponseWriter, r *http.Request) {
	println("ListIndustries")
	query := map[string]interface{}{}

	industrys, err := c.store.ListIndustries(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, industrys)
}

func (c *Controllers) CreateIndustry(w http.ResponseWriter, r *http.Request) {
	var industryPayload models.Industry

	errDecode := json.NewDecoder(r.Body).Decode(&industryPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	industry, err := c.store.CreateIndustry(industryPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, industry)
}

func (c *Controllers) ModifyIndustry(w http.ResponseWriter, r *http.Request) {
	var industryPayload models.Industry

	industryID := uuid.Must(uuid.Parse(chi.URLParam(r, "industryID")))

	errDecode := json.NewDecoder(r.Body).Decode(&industryPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	industry, err := c.store.ModifyIndustry(industryID, industryPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, industry)
}

func (c *Controllers) DeleteIndustry(w http.ResponseWriter, r *http.Request) {
	industryID := uuid.Must(uuid.Parse(chi.URLParam(r, "industryID")))

	err := c.store.DeleteIndustry(industryID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
