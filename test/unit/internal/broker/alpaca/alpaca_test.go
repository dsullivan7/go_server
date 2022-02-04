package zap_test

import (
	"go_server/internal/broker/alpaca"
	mockHTTP "go_server/test/mocks/http"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAlpacaCreateAccount(t *testing.T) {
	t.Parallel()

	mockHTTPClient := mockHTTP.NewClient()

	alpacaClient := alpaca.NewBroker(
		"alpacaAPIKey",
		"alpacaAPISecret",
		"alpacaAPIURL",
		mockHTTPClient,
	)

	mockHTTPClient.On("Do", mock.Anything).Return(&http.Response{StatusCode: 200}, nil)

	accountNumber, errAcc := alpacaClient.CreateAccount(
		"emailAddress",
		"phoneNumber",
	)

	assert.Nil(t, errAcc)
	assert.NotEqual(t, accountNumber, "")

	mockHTTPClient.AssertExpectations(t)
}
