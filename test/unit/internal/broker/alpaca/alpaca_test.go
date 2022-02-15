package alpaca_test

import (
	"bytes"
	"encoding/json"
	"go_server/internal/broker/alpaca"
	mockHTTP "go_server/test/mocks/http"
	"io/ioutil"
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

	body := map[string]interface{}{
		"id": "test",
	}

	jsonBytes, errMarshal := json.Marshal(body)

	assert.Nil(t, errMarshal)

	mockHTTPClient.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))),
		},
		nil,
	)

	alpacaAccountID, errAcc := alpacaClient.CreateAccount(
		"givenName",
		"familyName",
		"dateOfBirth",
		"taxID",
		"emailAddress",
		"phoneNumber",
		"streetAddress",
		"city",
		"state",
		"postalCode",
		"fundingSource",
		"ipAddress",
	)

	assert.Nil(t, errAcc)
	assert.Equal(t, alpacaAccountID, "test")

	mockHTTPClient.AssertExpectations(t)
}

func TestAlpacaGetAccount(t *testing.T) {
	t.Parallel()

	mockHTTPClient := mockHTTP.NewClient()

	alpacaClient := alpaca.NewBroker(
		"alpacaAPIKey",
		"alpacaAPISecret",
		"alpacaAPIURL",
		mockHTTPClient,
	)

	body := map[string]interface{}{
		"id": "test",
		"cash": "123.45",
	}

	jsonBytes, errMarshal := json.Marshal(body)

	assert.Nil(t, errMarshal)

	mockHTTPClient.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))),
		},
		nil,
	)

	account, errAcc := alpacaClient.GetAccount("test")

	assert.Nil(t, errAcc)
	assert.Equal(t, account.AccountID, "test")
	assert.Equal(t, account.Cash, 123.45)

	mockHTTPClient.AssertExpectations(t)
}
