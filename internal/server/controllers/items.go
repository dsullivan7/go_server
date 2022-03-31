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

func (c *Controllers) GetItem(w http.ResponseWriter, r *http.Request) {
	itemID := uuid.Must(uuid.Parse(chi.URLParam(r, "item_id")))

	item, err := c.store.GetItem(itemID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, item)
}

func (c *Controllers) ListItems(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	invoiceID := r.URL.Query().Get("invoice_id")

	if invoiceID != "" {
		query["invoice_id"] = invoiceID
	}

	items, err := c.store.ListItems(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, items)
}

func (c *Controllers) CreateItem(w http.ResponseWriter, r *http.Request) {
	var itemPayload models.Item

	errDecode := json.NewDecoder(r.Body).Decode(&itemPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	item, err := c.store.CreateItem(itemPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, item)
}

func (c *Controllers) ModifyItem(w http.ResponseWriter, r *http.Request) {
	var itemPayload models.Item

	itemID := uuid.Must(uuid.Parse(chi.URLParam(r, "item_id")))

	errDecode := json.NewDecoder(r.Body).Decode(&itemPayload)

	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	item, err := c.store.ModifyItem(itemID, itemPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, item)
}

func (c *Controllers) DeleteItem(w http.ResponseWriter, r *http.Request) {
	itemID := uuid.Must(uuid.Parse(chi.URLParam(r, "item_id")))

	err := c.store.DeleteItem(itemID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.NoContent(w, r)
}
