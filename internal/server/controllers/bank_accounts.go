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

func (c *Controllers) GetBankAccount(w http.ResponseWriter, r *http.Request) {
	bankAccountID := uuid.Must(uuid.Parse(chi.URLParam(r, "bank_account_id")))

	bankAccount, err := c.store.GetBankAccount(bankAccountID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, bankAccount)
}

func (c *Controllers) ListBankAccounts(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	userID := r.URL.Query().Get("user_id")

	if userID != "" {
		query["user_id"] = userID
	}

	bankAccounts, err := c.store.ListBankAccounts(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, bankAccounts)
}

func (c *Controllers) CreateBankAccount(w http.ResponseWriter, r *http.Request) {
	var bankAccountReq map[string]string

	errDecode := json.NewDecoder(r.Body).Decode(&bankAccountReq)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	accessToken, errAccessToken := c.bank.GetAccessToken(bankAccountReq["public_token"])

	if errAccessToken != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errAccessToken})

		return
	}

	userID := uuid.Must(uuid.Parse(bankAccountReq["user_id"]))

	name, errAccount := c.bank.GetAccount(accessToken)

	if errAccount != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errAccount})

		return
	}

	bankAccountPayload := models.BankAccount{
		UserID:           &userID,
		Name:             &name,
		AccessToken: 			&accessToken,
	}

	bankAccount, err := c.store.CreateBankAccount(bankAccountPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, bankAccount)
}

func (c *Controllers) ModifyBankAccount(w http.ResponseWriter, r *http.Request) {
	var bankAccountPayload models.BankAccount

	bankAccountID := uuid.Must(uuid.Parse(chi.URLParam(r, "bank_account_id")))

	errDecode := json.NewDecoder(r.Body).Decode(&bankAccountPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	bankAccount, err := c.store.ModifyBankAccount(bankAccountID, bankAccountPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, bankAccount)
}

func (c *Controllers) DeleteBankAccount(w http.ResponseWriter, r *http.Request) {
	bankAccountID := uuid.Must(uuid.Parse(chi.URLParam(r, "bank_account_id")))

	err := c.store.DeleteBankAccount(bankAccountID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
