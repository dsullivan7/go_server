package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetBrokerageAccount(brokerageAccountID uuid.UUID) (*models.BrokerageAccount, error) {
	var brokerageAccount models.BrokerageAccount

	err := gormStore.database.First(&brokerageAccount, brokerageAccountID).Error
	if err != nil {
		return nil, err
	}

	return &brokerageAccount, nil
}

func (gormStore *Store) ListBrokerageAccounts(query map[string]interface{}) ([]models.BrokerageAccount, error) {
	var brokerageAccounts []models.BrokerageAccount

	err := gormStore.database.Where(query).Order("created_at desc").Find(&brokerageAccounts).Error
	if err != nil {
		return nil, err
	}

	return brokerageAccounts, nil
}

func (gormStore *Store) CreateBrokerageAccount(
	brokerageAccountPayload models.BrokerageAccount,
) (*models.BrokerageAccount, error) {
	brokerageAccount := brokerageAccountPayload

	err := gormStore.database.Create(&brokerageAccount).Error
	if err != nil {
		return nil, err
	}

	return &brokerageAccount, nil
}

func (gormStore *Store) ModifyBrokerageAccount(
	brokerageAccountID uuid.UUID,
	brokerageAccountPayload models.BrokerageAccount,
) (*models.BrokerageAccount, error) {
	var brokerageAccountFound models.BrokerageAccount

	errFind := gormStore.database.Where("brokerage_account_id = ?", brokerageAccountID).First(&brokerageAccountFound).Error

	if errFind != nil {
		return nil, errFind
	}

	errUpdate := gormStore.database.Save(&brokerageAccountFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &brokerageAccountFound, nil
}

func (gormStore *Store) DeleteBrokerageAccount(brokerageAccountID uuid.UUID) error {
	err := gormStore.database.Delete(&models.BrokerageAccount{}, brokerageAccountID).Error

	return err
}
