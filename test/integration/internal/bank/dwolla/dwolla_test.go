package dwolla_test

import (
	// "fmt".
	"go_server/internal/bank/dwolla"
	"go_server/internal/config"
	goServerZapLogger "go_server/internal/logger/zap"

	// "go_server/internal/models".
	"testing"

	// "github.com/google/uuid".

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
		cfg.DwollaWebhookURL,
		cfg.DwollaWebhookSecret,
		logger,
	)

	webhook, errWebhook := dwollaBank.CreateWebhook()

	// randomID := uuid.New()
	//
	// firstName := "firstName"
	// lastName := "lastName"
	// dateOfBirth := "1980-01-01"
	// ssn := "666-55-4321"
	// email := fmt.Sprint("dbsullivan23+", randomID.String(), "@gmail.com")
	// phoneNumber := "555-444-3322"
	// address1 := "42 Faux St"
	// city := "New York"
	// state := "NY"
	// postalCode := "10009"
	//
	// dwollaCustomerID := "bab9537c-610e-46cf-b60b-0f92c2578764"
	//
	// user := models.User{
	//   FirstName: &firstName,
	//   LastName: &lastName,
	//   DateOfBirth: &dateOfBirth,
	//   SSN: &ssn,
	//   Email: &email,
	//   PhoneNumber: &phoneNumber,
	//   Address1: &address1,
	//   City: &city,
	//   State: &state,
	//   PostalCode: &postalCode,
	//   DwollaCustomerID: &dwollaCustomerID,
	// }

	// dwollaUser, errAcc := dwollaBank.CreateCustomer(user)

	// assert.Nil(t, errAcc)

	// processorToken  := "processor-sandbox-e6b301a9-35ec-4b5e-a68a-46d99eaba5ad"
	// bankAccount, errBank := dwollaBank.CreateBankAccount(user, processorToken)
	//
	// println("bankAccount.DwollaFundingSourceID")
	// println(*bankAccount.DwollaFundingSourceID)
	//
	// assert.Nil(t, errBank)
	// assert.NotNil(t, bankAccount.DwollaFundingSourceID)

	assert.Nil(t, errWebhook)
	assert.NotNil(t, webhook)
}
