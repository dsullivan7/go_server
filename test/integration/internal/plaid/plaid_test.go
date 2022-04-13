package plaid_test

import (
	"go_server/internal/config"
	goServerZapLogger "go_server/internal/logger/zap"
	goServerPlaid "go_server/internal/plaid"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestAlpacaCreateAccount(t *testing.T) {
	t.Skip("No integration")
	t.Parallel()

	cfg, configError := config.NewConfig()

	assert.Nil(t, configError)

	zapLogger, errZap := zap.NewProduction()

	assert.Nil(t, errZap)

	logger := goServerZapLogger.NewLogger(zapLogger)

	plaidClient := goServerPlaid.NewClient(
		cfg.PlaidClientID,
		cfg.PlaidSecret,
		cfg.PlaidAPIURL,
		cfg.PlaidRedirectURI,
		logger,
	)

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
