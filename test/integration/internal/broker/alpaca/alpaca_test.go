package zap_test

import (
	// "fmt"
	"go_server/internal/broker/alpaca"
	goServerHTTP "go_server/internal/http"
	"go_server/internal/config"
	"testing"

	// "github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestAlpacaCreateAccount(t *testing.T) {
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

	// randomID := uuid.New()
	//
	// accountNumber, errAcc := alpacaClient.CreateAccount(
	// 	"Blah",
	// 	"Blaherson",
	// 	"1990-01-01",
	// 	"666-55-4321",
	// 	fmt.Sprint("dbsullivan23+", randomID.String(), "@gmail.com"),
	// 	"555-444-3322",
	// 	"42 Faux St",
	// 	"New York",
	// 	"NY",
	// 	"10009",
	// 	"savings",
	// 	"185.13.21.99",
	// )
	//
	// assert.Nil(t, errAcc)
	// assert.NotEqual(t, accountNumber, "")

	// accountNumber := alpacaClient.DeleteAccount(accountNumber)
	// errDel := alpacaClient.DeleteAccount(accountNumber)
	// assert.Nil(t, errDel)

	accountNumber := "a387ec8a-f4f4-430b-a619-48c4f81222d2"
	orderID, errOrder := alpacaClient.CreateOrder(
		accountNumber,
		"AAPL",
		1.23,
		"buy",
	)
	println("orderID")
	println(orderID)
	assert.Nil(t, errOrder)
}
