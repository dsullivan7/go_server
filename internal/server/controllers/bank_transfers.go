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

func (c *Controllers) GetBankTransfer(w http.ResponseWriter, r *http.Request) {
	bankTransferID := uuid.Must(uuid.Parse(chi.URLParam(r, "bank_transfer_id")))

	bankTransfer, err := c.store.GetBankTransfer(bankTransferID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, bankTransfer)
}

func (c *Controllers) ListBankTransfers(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	userID := r.URL.Query().Get("user_id")

	if userID != "" {
		query["user_id"] = userID
	}

	bankTransfers, err := c.store.ListBankTransfers(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, bankTransfers)
}

func (c *Controllers) CreateBankTransfer(w http.ResponseWriter, r *http.Request) {
	var bankTransferPayload models.BankTransfer

	errDecode := json.NewDecoder(r.Body).Decode(&bankTransferPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	bankTransfer, err := c.store.CreateBankTransfer(bankTransferPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, bankTransfer)
}

func (c *Controllers) ModifyBankTransfer(w http.ResponseWriter, r *http.Request) {
	var bankTransferPayload models.BankTransfer

	bankTransferID := uuid.Must(uuid.Parse(chi.URLParam(r, "bank_transfer_id")))

	errDecode := json.NewDecoder(r.Body).Decode(&bankTransferPayload)

	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	bankTransfer, err := c.store.ModifyBankTransfer(bankTransferID, bankTransferPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, bankTransfer)
}

func (c *Controllers) DeleteBankTransfer(w http.ResponseWriter, r *http.Request) {
	bankTransferID := uuid.Must(uuid.Parse(chi.URLParam(r, "bank_transfer_id")))

	err := c.store.DeleteBankTransfer(bankTransferID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.NoContent(w, r)
}
