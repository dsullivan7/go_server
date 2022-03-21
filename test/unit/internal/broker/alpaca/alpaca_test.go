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

	body := [2]map[string]interface{}{
		{
			"asset_id":     "asset_id_1",
			"market_value": "345.67",
			"symbol":       "symbol_1",
		},
		{
			"asset_id":     "asset_id_2",
			"market_value": "123.45",
			"symbol":       "symbol_2",
		},
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

	positions, errPositions := alpacaClient.ListPositions("test")

	assert.Nil(t, errPositions)
	assert.Equal(t, len(positions), 2)

	assert.Equal(t, positions[0].PositionID, "asset_id_1")
	assert.Equal(t, positions[0].Symbol, "symbol_1")
	assert.Equal(t, positions[0].MarketValue, 345.67)

	assert.Equal(t, positions[1].PositionID, "asset_id_2")
	assert.Equal(t, positions[1].Symbol, "symbol_2")
	assert.Equal(t, positions[1].MarketValue, 123.45)

	mockHTTPClient.AssertExpectations(t)
}
