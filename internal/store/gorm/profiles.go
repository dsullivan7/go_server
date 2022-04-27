package gorm

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

func (gormStore *Store) GetProfile(profileID uuid.UUID) (*models.Profile, error) {
	var profile models.Profile

	err := gormStore.database.First(&profile, profileID).Error
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (gormStore *Store) ListProfiles(query map[string]interface{}) ([]models.Profile, error) {
	var profiles []models.Profile

	err := gormStore.database.Where(query).Order("created_at desc").Find(&profiles).Error
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (gormStore *Store) CreateProfile(profilePayload models.Profile) (*models.Profile, error) {
	profile := profilePayload

	err := gormStore.database.Create(&profile).Error
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (gormStore *Store) ModifyProfile(
	profileID uuid.UUID,
	profilePayload models.Profile,
) (*models.Profile, error) {
	var profileFound models.Profile

	errFind := gormStore.database.Where("profile_id = ?", profileID).First(&profileFound).Error

	if errFind != nil {
		return nil, errFind
	}

	if profilePayload.Username != "" {
		profileFound.Username = profilePayload.Username
	}

	errUpdate := gormStore.database.Save(&profileFound).Error

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &profileFound, nil
}

func (gormStore *Store) DeleteProfile(profileID uuid.UUID) error {
	err := gormStore.database.Delete(&models.Profile{}, profileID).Error

	return err
}
