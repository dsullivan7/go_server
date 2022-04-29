package controllers

import (
	"encoding/json"
	"go_server/internal/errors"
	"go_server/internal/models"
	"net/http"
)

func (c *Controllers) DwollaWebhook(w http.ResponseWriter, r *http.Request) {
	var webhookPayload map[string]interface{}

	errDecode := json.NewDecoder(r.Body).Decode(&webhookPayload)

	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	if webhookPayload["topic"] == "customer_transfer_completed" {
		dwollaTransferID := webhookPayload["resourceId"]
		transfers, errList := c.store.ListBankTransfers(map[string]interface{}{"dwolla_transfer_id": dwollaTransferID})

		if errList != nil {
			c.utils.HandleError(w, r, errors.HTTPUserError{Err: errList})

			return
		}

		transfer := transfers[0]
		_, errModify := c.store.ModifyBankTransfer(transfer.BankTransferID, models.BankTransfer{Status: "complete"})

		if errModify != nil {
			c.utils.HandleError(w, r, errors.HTTPUserError{Err: errModify})

			return
		}
	}
}
