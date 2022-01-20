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
	backAccountID := uuid.Must(uuid.Parse(chi.URLParam(r, "back_account_id")))

	backAccount, err := c.store.GetBankAccount(backAccountID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, backAccount)
}

func (c *Controllers) ListBankAccounts(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	portfolioId := r.URL.Query().Get("user_id")

	if portfolioId != "" {
		query["user_id"] = portfolioId
	}

	backAccounts, err := c.store.ListBankAccounts(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, backAccounts)
}

func (c *Controllers) CreateBankAccount(w http.ResponseWriter, r *http.Request) {
	var backAccountPayload models.BankAccount

	errDecode := json.NewDecoder(r.Body).Decode(&backAccountPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	backAccount, err := c.store.CreateBankAccount(backAccountPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, backAccount)
}

func (c *Controllers) ModifyBankAccount(w http.ResponseWriter, r *http.Request) {
	var backAccountPayload models.BankAccount

	backAccountID := uuid.Must(uuid.Parse(chi.URLParam(r, "back_account_id")))

	errDecode := json.NewDecoder(r.Body).Decode(&backAccountPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	backAccount, err := c.store.ModifyBankAccount(backAccountID, backAccountPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, backAccount)
}

func (c *Controllers) DeleteBankAccount(w http.ResponseWriter, r *http.Request) {
	backAccountID := uuid.Must(uuid.Parse(chi.URLParam(r, "back_account_id")))

	err := c.store.DeleteBankAccount(backAccountID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
