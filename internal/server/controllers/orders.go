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

func (c *Controllers) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := uuid.Must(uuid.Parse(chi.URLParam(r, "order_id")))

	order, err := c.store.GetOrder(orderID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, order)
}

func (c *Controllers) ListOrders(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}

	orders, err := c.store.ListOrders(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, orders)
}

func (c *Controllers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderPayload models.Order

	errDecode := json.NewDecoder(r.Body).Decode(&orderPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	order, err := c.store.CreateOrder(orderPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, order)
}

func (c *Controllers) ModifyOrder(w http.ResponseWriter, r *http.Request) {
	var orderPayload models.Order

	orderID := uuid.Must(uuid.Parse(chi.URLParam(r, "order_id")))

	errDecode := json.NewDecoder(r.Body).Decode(&orderPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	order, err := c.store.ModifyOrder(orderID, orderPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, order)
}

func (c *Controllers) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderID := uuid.Must(uuid.Parse(chi.URLParam(r, "order_id")))

	err := c.store.DeleteOrder(orderID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.NoContent(w, r)
}
