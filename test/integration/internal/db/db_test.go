package db_test

import (
	"go_server/internal/config"
	"go_server/internal/db"
	"go_server/internal/models"
	"go_server/test/utils"
	"testing"

	"github.com/stretchr/testify/assert"
  "github.com/google/uuid"
)

func TestDBIntegration(parentT *testing.T) {
  parentT.Parallel()

  cfg, errConfig := config.NewConfig()

  assert.Nil(parentT, errConfig)

  connection, errConnection := db.NewSQLConnection(
    cfg.DBHost,
    cfg.DBName,
    cfg.DBPort,
    cfg.DBUser,
    cfg.DBPassword,
    cfg.DBSSL,
  )
  assert.Nil(parentT, errConnection)

  db, errDB := db.NewGormDB(connection)
  assert.Nil(parentT, errDB)

  dbUtility := utils.NewSQLDatabaseUtility(connection)

  errTruncateBefore := dbUtility.TruncateAll()
  assert.Nil(parentT, errTruncateBefore)

  defer func() {
    errTruncateAfter := dbUtility.TruncateAll()
    assert.Nil(parentT, errTruncateAfter)
  }()

  parentT.Run("User", func(t *testing.T) {
    t.Parallel()
    firstName := "firstName"
    lastName := "lastName"
    auth0ID := uuid.New().String()

    user := models.User{
      FirstName: &firstName,
      LastName: &lastName,
      Auth0ID: &auth0ID,
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
