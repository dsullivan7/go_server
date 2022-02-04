package zap_test

import (
	"go_server/internal/broker/alpaca"
	goServerHTTP "go_server/internal/http"
	"go_server/internal/config"
	"testing"

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

	accountNumber, errAcc := alpacaClient.CreateAccount(
		"emailAddress",
		"phoneNumber",
	)

	assert.Nil(t, errAcc)
	assert.NotEqual(t, accountNumber, "")
}
