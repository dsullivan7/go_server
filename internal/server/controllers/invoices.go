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

func (c *Controllers) GetInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceID := uuid.Must(uuid.Parse(chi.URLParam(r, "invoice_id")))

	invoice, err := c.store.GetInvoice(invoiceID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, invoice)
}

func (c *Controllers) ListInvoices(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	userID := r.URL.Query().Get("user_id")

	if userID != "" {
		query["user_id"] = userID
	}

	invoices, err := c.store.ListInvoices(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, invoices)
}

func (c *Controllers) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var invoicePayload models.Invoice

	errDecode := json.NewDecoder(r.Body).Decode(&invoicePayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	invoice, err := c.store.CreateInvoice(invoicePayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, invoice)
}

func (c *Controllers) ModifyInvoice(w http.ResponseWriter, r *http.Request) {
	var invoicePayload models.Invoice

	invoiceID := uuid.Must(uuid.Parse(chi.URLParam(r, "invoice_id")))

	errDecode := json.NewDecoder(r.Body).Decode(&invoicePayload)

	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	invoice, err := c.store.ModifyInvoice(invoiceID, invoicePayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, invoice)
}

func (c *Controllers) DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceID := uuid.Must(uuid.Parse(chi.URLParam(r, "invoice_id")))

	err := c.store.DeleteInvoice(invoiceID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.NoContent(w, r)
}
