package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetInvoice(invoiceID uuid.UUID) (*models.Invoice, error) {
	var invoice models.Invoice

	err := gormStore.database.First(&invoice, invoiceID).Error
	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (gormStore *Store) ListInvoices(query map[string]interface{}) ([]models.Invoice, error) {
	var invoices []models.Invoice

	err := gormStore.database.Where(query).Order("created_at desc").Find(&invoices).Error
	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (gormStore *Store) CreateInvoice(invoicePayload models.Invoice) (*models.Invoice, error) {
	invoice := invoicePayload

	err := gormStore.database.Create(&invoice).Error
	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (gormStore *Store) ModifyInvoice(
	invoiceID uuid.UUID,
	invoicePayload models.Invoice,
) (*models.Invoice, error) {
	var invoiceFound models.Invoice

	errFind := gormStore.database.Where("invoice_id = ?", invoiceID).First(&invoiceFound).Error

	if errFind != nil {
		return nil, errFind
	}

	if invoicePayload.BankAccountID != nil {
		invoiceFound.BankAccountID = invoicePayload.BankAccountID
	}

	errUpdate := gormStore.database.Save(&invoiceFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &invoiceFound, nil
}

func (gormStore *Store) DeleteInvoice(invoiceID uuid.UUID) error {
	err := gormStore.database.Delete(&models.Invoice{}, invoiceID).Error

	return err
}
