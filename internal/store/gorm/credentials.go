package gorm

import (
	"go_server/internal/models"
)

func (gormStore *Store) CreateCredential(credentialPayload models.Credential) (*models.Credential, error) {
	credential := credentialPayload

	err := gormStore.database.Create(&credential).Error
	if err != nil {
		return nil, err
	}

	return &credential, nil
}
