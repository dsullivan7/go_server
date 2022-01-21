package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetBankAccount(bankAccountID uuid.UUID) (*models.BankAccount, error) {
	var bankAccount models.BankAccount

	err := gormStore.database.First(&bankAccount, bankAccountID).Error
	if err != nil {
		return nil, err
	}

	return &bankAccount, nil
}

func (gormStore *Store) ListBankAccounts(query map[string]interface{}) ([]models.BankAccount, error) {
	var bankAccounts []models.BankAccount

	err := gormStore.database.Where(query).Order("created_at desc").Find(&bankAccounts).Error
	if err != nil {
		return nil, err
	}

	return bankAccounts, nil
}

func (gormStore *Store) CreateBankAccount(bankAccountPayload models.BankAccount) (*models.BankAccount, error) {
	bankAccount := bankAccountPayload

	err := gormStore.database.Create(&bankAccount).Error
	if err != nil {
		return nil, err
	}

	return &bankAccount, nil
}

func (gormStore *Store) ModifyBankAccount(
	bankAccountID uuid.UUID,
	bankAccountPayload models.BankAccount,
) (*models.BankAccount, error) {
	var bankAccountFound models.BankAccount

	errFind := gormStore.database.Where("bank_account_id = ?", bankAccountID).First(&bankAccountFound).Error

	if errFind != nil {
		return nil, errFind
	}

	if bankAccountPayload.UserID != nil {
		bankAccountFound.UserID = bankAccountPayload.UserID
	}

	errUpdate := gormStore.database.Save(&bankAccountFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &bankAccountFound, nil
}

func (gormStore *Store) DeleteBankAccount(bankAccountID uuid.UUID) error {
	err := gormStore.database.Delete(&models.BankAccount{}, bankAccountID).Error

	return err
}
