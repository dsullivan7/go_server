package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetBankTransfer(bankTransferID uuid.UUID) (*models.BankTransfer, error) {
	var bankTransfer models.BankTransfer

	err := gormStore.database.First(&bankTransfer, bankTransferID).Error
	if err != nil {
		return nil, err
	}

	return &bankTransfer, nil
}

func (gormStore *Store) ListBankTransfers(query map[string]interface{}) ([]models.BankTransfer, error) {
	var bankTransfers []models.BankTransfer

	err := gormStore.database.Where(query).Order("created_at desc").Find(&bankTransfers).Error
	if err != nil {
		return nil, err
	}

	return bankTransfers, nil
}

func (gormStore *Store) CreateBankTransfer(bankTransferPayload models.BankTransfer) (*models.BankTransfer, error) {
	bankTransfer := bankTransferPayload

	err := gormStore.database.Create(&bankTransfer).Error
	if err != nil {
		return nil, err
	}

	return &bankTransfer, nil
}

func (gormStore *Store) ModifyBankTransfer(
	bankTransferID uuid.UUID,
	bankTransferPayload models.BankTransfer,
) (*models.BankTransfer, error) {
	var bankTransferFound models.BankTransfer

	errFind := gormStore.database.Where("bank_account_id = ?", bankTransferID).First(&bankTransferFound).Error

	if errFind != nil {
		return nil, errFind
	}

	if bankTransferPayload.UserID != nil {
		bankTransferFound.UserID = bankTransferPayload.UserID
	}

	errUpdate := gormStore.database.Save(&bankTransferFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &bankTransferFound, nil
}

func (gormStore *Store) DeleteBankTransfer(bankTransferID uuid.UUID) error {
	err := gormStore.database.Delete(&models.BankTransfer{}, bankTransferID).Error

	return err
}
