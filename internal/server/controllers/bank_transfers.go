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
	var bankTransferReq map[string]interface{}

	errDecode := json.NewDecoder(r.Body).Decode(&bankTransferReq)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	alpacaTransferID, errTransfer := c.broker.CreateTransfer(
		bankTransferReq["alpaca_account_id"].(string),
		bankTransferReq["alpaca_ach_relationship_id"].(string),
		bankTransferReq["amount"].(float64),
		"INCOMING",
	)

	if errTransfer != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errTransfer})

		return
	}

	userID := uuid.Must(uuid.Parse(bankTransferReq["user_id"].(string)))

  bankTransferPayload := models.BankTransfer{
    UserID: &userID,
    Amount: bankTransferReq["amount"].(float64),
    Status: "PENDING",
    AlpacaTransferID: &alpacaTransferID,
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
