package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (mockStore *MockStore) GetInvoice(invoiceID uuid.UUID) (*models.Invoice, error) {
	args := mockStore.Called(invoiceID)

	return args.Get(0).(*models.Invoice), args.Error(1)
}

func (mockStore *MockStore) ListInvoices(query map[string]interface{}) ([]models.Invoice, error) {
	args := mockStore.Called(query)

	return args.Get(0).([]models.Invoice), args.Error(1)
}

func (mockStore *MockStore) CreateInvoice(invoicePayload models.Invoice) (*models.Invoice, error) {
	args := mockStore.Called(invoicePayload)

	return args.Get(0).(*models.Invoice), args.Error(1)
}

func (mockStore *MockStore) ModifyInvoice(
	invoiceID uuid.UUID,
	invoicePayload models.Invoice,
) (*models.Invoice, error) {
	args := mockStore.Called(invoiceID, invoicePayload)

	return args.Get(0).(*models.Invoice), args.Error(1)
}

func (mockStore *MockStore) DeleteInvoice(invoiceID uuid.UUID) error {
	args := mockStore.Called(invoiceID)

	return args.Error(0)
}
