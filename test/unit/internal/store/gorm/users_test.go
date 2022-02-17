package gorm_test

import (
	"time"
	"regexp"
	"go_server/internal/models"
	goServerGormStore "go_server/internal/store/gorm"
  "github.com/DATA-DOG/go-sqlmock"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
	"testing"

	"github.com/stretchr/testify/assert"
  "github.com/google/uuid"
)

func TestDBIntegration(parentT *testing.T) {
  parentT.Parallel()

  sqlDB, mock, errSQLMock := sqlmock.New()
  assert.Nil(parentT, errSQLMock)

  dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	})

  db, errDB := gorm.Open(dialector, &gorm.Config{})
  assert.Nil(parentT, errDB)

  store := goServerGormStore.NewStore(db)

  parentT.Run("User", func(t *testing.T) {
    t.Parallel()

    userID := uuid.New()
    firstName := "firstName"
    lastName := "lastName"
    auth0ID := uuid.New().String()
    createdAt := time.Now()
    updatedAt := time.Now()

    mock.ExpectBegin()
    mock.ExpectQuery(
      regexp.QuoteMeta(`
        INSERT INTO "users" ("auth0_id","first_name","last_name")
        VALUES ($1,$2,$3)
        RETURNING "user_id","created_at","updated_at"
       `)).
      WithArgs(auth0ID, firstName, lastName).
      WillReturnRows(
        sqlmock.NewRows([]string{"user_id", "created_at", "updated_at"}).
          AddRow(userID, createdAt, updatedAt))
    mock.ExpectCommit()

    // mock.ExpectExec("INSERT INTO users").
    //   WithArgs(firstName, lastName, auth0ID).
    //   WillReturnResult(sqlmock.NewResult(1, 1))

    user := models.User{
      FirstName: &firstName,
      LastName: &lastName,
      Auth0ID: &auth0ID,
    }

    userCreated, err := store.CreateUser(user)
    assert.Nil(t, err)

    assert.Equal(t, userCreated.UserID, userID)
    assert.Equal(t, *userCreated.FirstName, firstName)
    assert.Equal(t, *userCreated.LastName, lastName)
    assert.Equal(t, *userCreated.Auth0ID, auth0ID)
    assert.WithinDuration(t, userCreated.CreatedAt, createdAt, 0)
    assert.WithinDuration(t, userCreated.UpdatedAt, updatedAt, 0)

    errExpectations := mock.ExpectationsWereMet()
    assert.Nil(t, errExpectations)
  })
}
