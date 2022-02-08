package plaid_test

import (
	"go_server/internal/config"
	goServerPlaid "go_server/internal/plaid"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/plaid/plaid-go/plaid"
)

func TestAlpacaCreateAccount(t *testing.T) {
	t.Skip("No integration")
	t.Parallel()

	cfg, configError := config.NewConfig()

	assert.Nil(t, configError)

	plaidConfig := plaid.NewConfiguration()
	plaidConfig.AddDefaultHeader("PLAID-CLIENT-ID", cfg.PlaidClientID)
	plaidConfig.AddDefaultHeader("PLAID-SECRET", cfg.PlaidSecret)
	plaidConfig.UseEnvironment(plaid.Sandbox)
	plaidAPIClient := plaid.NewAPIClient(plaidConfig)
	plaidClient := goServerPlaid.NewClient(plaidAPIClient, cfg.PlaidRedirectURI)

	accountID := "NoLzzvE5x8sKzLzWk86eHXARMRKB43fWnLrDz"
	accessToken := "access-sandbox-52c18719-eafa-4198-9ef9-f32fdd3300f6"

	processorToken, errRel := plaidClient.CreateProcessorToken(
		accessToken,
		accountID,
		"alpaca",
	)

	assert.NotNil(t, processorToken)
	assert.Nil(t, errRel)
}
