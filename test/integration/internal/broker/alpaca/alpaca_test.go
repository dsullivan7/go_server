package alpaca_test

import (
	"fmt"
	"go_server/internal/broker/alpaca"
	"go_server/internal/config"
	goServerHTTP "go_server/internal/http"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestAlpacaCreateAccount(t *testing.T) {
	t.Skip("No integration")
	t.Parallel()

	cfg, configError := config.NewConfig()

	assert.Nil(t, configError)

	httpClient := goServerHTTP.NewClient()

	alpacaClient := alpaca.NewBroker(
		cfg.AlpacaAPIKey,
		cfg.AlpacaAPISecret,
		cfg.AlpacaAPIURL,
		httpClient,
	)

	randomID := uuid.New()

	accountNumber, errAcc := alpacaClient.CreateAccount(
		"Blah",
		"Blaherson",
		"1990-01-01",
		"666-55-4321",
		fmt.Sprint("dbsullivan23+", randomID.String(), "@gmail.com"),
		"555-444-3322",
		"42 Faux St",
		"New York",
		"NY",
		"10009",
		"savings",
		"185.13.21.99",
	)

	assert.Nil(t, errAcc)
	assert.NotEqual(t, accountNumber, "")

	errDel := alpacaClient.DeleteAccount(accountNumber)
	assert.Nil(t, errDel)

	// accountNumber := "a387ec8a-f4f4-430b-a619-48c4f81222d2"

	relactionshipID, errRel := alpacaClient.CreateACHRelationship(
		accountNumber,
		"processor-sandbox-c9f270e8-1d76-4f98-b6b8-c6b81a2708c6",
	)
	assert.Nil(t, errRel)

	// relactionshipID := "7b022285-d2fe-4c06-9fa6-f2d7df2906eb"
	transferID, errTransfer := alpacaClient.CreateTransfer(
		accountNumber,
		relactionshipID,
		5000,
		"INCOMING",
	)
	assert.Nil(t, errTransfer)
	assert.NotNil(t, transferID)

	// transferID := "1506d87a-53cb-46ea-848f-b0a9a6059a19"

	orderID, errOrder := alpacaClient.CreateOrder(
		accountNumber,
		"AAPL",
		1.23,
		"buy",
	)
	assert.Nil(t, errOrder)
	assert.NotNil(t, orderID)
}
