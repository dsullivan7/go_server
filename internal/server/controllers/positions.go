package controllers

import (
	"go_server/internal/errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/go-chi/render"
)

func (c *Controllers) ListPositions(w http.ResponseWriter, r *http.Request) {
	brokerageAccountID := uuid.Must(uuid.Parse(r.URL.Query().Get("brokerage_account_id")))

	brokerageAccount, errBrokerageAccount := c.store.GetBrokerageAccount(brokerageAccountID)

	if errBrokerageAccount != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errBrokerageAccount})

		return
	}

	positions, errPositions := c.broker.ListPositions(*brokerageAccount.AlpacaAccountID)

	if errPositions != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errPositions})

		return
	}

	render.JSON(w, r, positions)
}
