package controllers

import (
	"encoding/json"
	"go_server/internal/errors"
	"go_server/internal/models"
	"net/http"
	"time"

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
	userID := r.URL.Query().Get("user_id")

	if userID != "" {
		query["user_id"] = userID
	}

	orders, err := c.store.ListOrders(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, orders)
}

func (c *Controllers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq map[string]interface{}

	errDecode := json.NewDecoder(r.Body).Decode(&orderReq)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	userID := uuid.Must(uuid.Parse(orderReq["user_id"].(string)))

	completedAt := time.Now()

	orderPayload := models.Order{
		UserID:      &userID,
		Amount:      int(orderReq["amount"].(float64)),
		Side:        orderReq["side"].(string),
		Status:      "complete",
		CompletedAt: &completedAt,
	}

	order, errOrder := c.store.CreateOrder(orderPayload)

	if errOrder != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errOrder})

		return
	}

	childOrderPayload := models.Order{
		ParentOrderID: &order.OrderID,
		UserID:        &userID,
		Amount:        int(orderReq["amount"].(float64)),
		Side:          orderReq["side"].(string),
		Status:        "complete",
		CompletedAt:   &completedAt,
	}

	_, errChildOrder := c.store.CreateOrder(childOrderPayload)

	if errChildOrder != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errChildOrder})

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
