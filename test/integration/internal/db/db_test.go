package db_test

import (
	"go_server/internal/config"
	"go_server/internal/db"
	"go_server/internal/models"
	"go_server/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDBIntegration(tParent *testing.T) {
	tParent.Parallel()

	cfg, errConfig := config.NewConfig()

	assert.Nil(tParent, errConfig)

	connection, errConnection := db.NewSQLConnection(
		cfg.DBHost,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBSSL,
	)
	assert.Nil(tParent, errConnection)

	db, errDB := db.NewGormDB(connection)
	assert.Nil(tParent, errDB)

	dbUtility := utils.NewSQLDatabaseUtility(connection)

	errTruncateBefore := dbUtility.TruncateAll()
	assert.Nil(tParent, errTruncateBefore)

	tParent.Cleanup(func() {
		errTruncateAfter := dbUtility.TruncateAll()
		assert.Nil(tParent, errTruncateAfter)
	})

	tParent.Run("User", func(t *testing.T) {
		t.Parallel()
		firstName := "firstName"
		lastName := "lastName"
		auth0ID := uuid.New().String()

		user := models.User{
			FirstName: &firstName,
			LastName:  &lastName,
			Auth0ID:   &auth0ID,
		}

		err := db.Create(&user).Error
		assert.Nil(t, err)

		assert.Equal(t, *user.FirstName, firstName)
		assert.Equal(t, *user.LastName, lastName)
		assert.Equal(t, *user.Auth0ID, auth0ID)
		assert.NotNil(t, user.CreatedAt)
		assert.NotNil(t, user.UpdatedAt)
		assert.NotNil(t, user.UserID)
	})

	tParent.Run("Order", func(t *testing.T) {
		t.Parallel()
		firstName := "firstName"
		lastName := "lastName"
		auth0ID := uuid.New().String()

		user := models.User{
			FirstName: &firstName,
			LastName:  &lastName,
			Auth0ID:   &auth0ID,
		}

		err := db.Create(&user).Error
		assert.Nil(t, err)

		assert.Equal(t, *user.FirstName, firstName)
		assert.Equal(t, *user.LastName, lastName)
		assert.Equal(t, *user.Auth0ID, auth0ID)
		assert.NotNil(t, user.CreatedAt)
		assert.NotNil(t, user.UpdatedAt)
		assert.NotNil(t, user.UserID)
	})
}
