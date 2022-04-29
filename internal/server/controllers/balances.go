package controllers

import (
	"go_server/internal/errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

const interestRate = 0.05

func (c *Controllers) GetBalances(w http.ResponseWriter, r *http.Request) {
	userID := uuid.Must(uuid.Parse(chi.URLParam(r, "user_id")))

	orders, err := c.store.ListChildOrders(userID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	total, principal, interest := c.services.GetBalance(orders, interestRate, time.Now())

	response := map[string]int{
		"total":     total,
		"principal": principal,
		"interest":  interest,
	}

	render.JSON(w, r, response)
}
