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

	brokerageAccountID := uuid.Must(uuid.Parse(bankAccountReq["brokerage_account_id"]))

	brokerageAccount, errBrokerageAccount := c.store.GetBrokerageAccount(brokerageAccountID)

	if errBrokerageAccount != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errBrokerageAccount})

		return
	}

	plaidAccessToken, errAccessToken := c.plaidClient.GetAccessToken(bankAccountReq["plaid_public_token"])

	if errAccessToken != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errAccessToken})

		return
	}

	plaidAccountID, name, errAccount := c.plaidClient.GetAccount(plaidAccessToken)

	if errAccount != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errAccount})

		return
	}

	plaidProcessorToken, errPlaidProcessorToken := c.plaidClient.CreateProcessorToken(plaidAccessToken, plaidAccountID, "alpaca")

	if errPlaidProcessorToken != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errPlaidProcessorToken})

		return
	}

	alpacaACHRelationshipID, errACHRelationship := c.broker.CreateACHRelationship(*brokerageAccount.AlpacaAccountID, plaidProcessorToken)

	if errACHRelationship != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errACHRelationship})

		return
	}

	userID := uuid.Must(uuid.Parse(bankAccountReq["user_id"]))

	bankAccountPayload := models.BankAccount{
		UserID:           &userID,
		Name:             &name,
		PlaidAccountID:   &plaidAccountID,
		PlaidAccessToken: &plaidAccessToken,
		AlpacaACHRelationshipID: &alpacaACHRelationshipID,
	}

	bankAccount, err := c.store.CreateBankAccount(bankAccountPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.Status(r, http.StatusCreated)
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

	render.NoContent(w, r)
}
