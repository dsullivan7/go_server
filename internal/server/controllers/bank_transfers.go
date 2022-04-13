package controllers

import (
	// "fmt".
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

	userID := uuid.Must(uuid.Parse(bankTransferReq["user_id"].(string)))

	// user, errUser := c.store.GetUser(userID)
	//
	// if errUser != nil {
	// 	c.utils.HandleError(w, r, errors.HTTPUserError{Err: errUser})
	//
	// 	return
	// }
	//
	// var plaidOriginationAccountID string
	// if (bankTransferReq["plaid_origination_account_id"] != nil) {
	// 	val, ok := bankTransferReq["plaid_origination_account_id"].(string)
	// 	if (ok) {
	// 		plaidOriginationAccountID = val
	// 	}
	// }
	//
	// plaidTransferAuthorizationID, errTransferAuth := c.plaidClient.CreateTransferAuthorization(
	// 	bankTransferReq["plaid_account_id"].(string),
	// 	bankTransferReq["plaid_access_token"].(string),
	// 	plaidOriginationAccountID,
	// 	bankTransferReq["amount"].(string),
	// 	"debit",
	// 	fmt.Sprint(user.FirstName, " ", user.LastName),
	// )
	//
	// if errTransferAuth != nil {
	// 	c.utils.HandleError(w, r, errors.HTTPUserError{Err: errTransferAuth})
	//
	// 	return
	// }
	//
	// plaidTransferID, errTransfer := c.plaidClient.CreateTransfer(
	// 	bankTransferReq["plaid_account_id"].(string),
	// 	bankTransferReq["plaid_access_token"].(string),
	// 	plaidOriginationAccountID,
	// 	plaidTransferAuthorizationID,
	// 	bankTransferReq["amount"].(string),
	// 	"debit",
	// 	fmt.Sprint(user.FirstName, " ", user.LastName),
	// )
	//
	// if errTransfer != nil {
	// 	c.utils.HandleError(w, r, errors.HTTPUserError{Err: errTransfer})
	//
	// 	return
	// }
	//
	var plaidTransferID string

	bankTransferPayload := models.BankTransfer{
		UserID:          &userID,
		Amount:          int(bankTransferReq["amount"].(float64)),
		Status:          "PENDING",
		PlaidTransferID: &plaidTransferID,
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
