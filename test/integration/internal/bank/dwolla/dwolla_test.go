package dwolla_test

import (
	"fmt"
  goServerZapLogger "go_server/internal/logger/zap"
	"go_server/internal/bank/dwolla"
	"go_server/internal/config"
	"go_server/internal/models"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
  "go.uber.org/zap"
)

func TestDwollaCreateAccount(t *testing.T) {
	t.Skip("No integration")
	t.Parallel()

	cfg, configError := config.NewConfig()

	assert.Nil(t, configError)

  zapLogger, errZap := zap.NewProduction()

	assert.Nil(t, errZap)

	logger := goServerZapLogger.NewLogger(zapLogger)

	dwollaBank := dwolla.NewBank(
		cfg.DwollaAPIKey,
		cfg.DwollaAPISecret,
		cfg.DwollaAPIURL,
    logger,
	)

	randomID := uuid.New()

  firstName := "firstName"
  lastName := "lastName"
  dateOfBirth := "1980-01-01"
  ssn := "666-55-4321"
  email := fmt.Sprint("dbsullivan23+", randomID.String(), "@gmail.com")
  phoneNumber := "555-444-3322"
  address1 := "42 Faux St"
  city := "New York"
  state := "NY"
  postalCode := "10009"

  dwollaCustomerID := "bab9537c-610e-46cf-b60b-0f92c2578764"

  user := models.User{
    FirstName: &firstName,
    LastName: &lastName,
    DateOfBirth: &dateOfBirth,
    SSN: &ssn,
    Email: &email,
    PhoneNumber: &phoneNumber,
    Address1: &address1,
    City: &city,
    State: &state,
    PostalCode: &postalCode,
    DwollaCustomerID: &dwollaCustomerID,
  }

	// dwollaUser, errAcc := dwollaBank.CreateCustomer(user)

	// assert.Nil(t, errAcc)

	bankAccount, errBank := dwollaBank.CreateBank(user, "")
	assert.Nil(t, errBank)
	assert.NotNil(t, bankAccount.DwollaFundingSourceID)
}
