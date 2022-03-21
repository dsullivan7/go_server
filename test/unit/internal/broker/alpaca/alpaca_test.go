package alpaca_test

import (
	"encoding/json"
	"fmt"
	"go_server/internal/broker/alpaca"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlpacaCreateAccount(t *testing.T) {
	t.Parallel()

	body := map[string]interface{}{
		"id": "test",
	}

	jsonBytes, errMarshal := json.Marshal(body)

	assert.Nil(t, errMarshal)

	expectedPath := "/v1/accounts"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, expectedPath)
		assert.Equal(t, r.Method, http.MethodPost)

		assert.NotNil(t, r.Header.Get("Authorization"))

		decoder := json.NewDecoder(r.Body)

		var alpacaResponse interface{}

		errDecode := decoder.Decode(&alpacaResponse)

		assert.Nil(t, errDecode)

		assert.Equal(
			t,
			alpacaResponse.(map[string]interface{})["contact"].(map[string]interface{})["email_address"],
			"emailAddress",
		)

		assert.Equal(
			t,
			alpacaResponse.(map[string]interface{})["contact"].(map[string]interface{})["phone_number"],
			"phoneNumber",
		)

		assert.Equal(
			t,
			alpacaResponse.(map[string]interface{})["contact"].(map[string]interface{})["street_address"].([]interface{})[0],
			"streetAddress",
		)

		assert.Equal(
			t,
			alpacaResponse.(map[string]interface{})["contact"].(map[string]interface{})["city"],
			"city",
		)

		assert.Equal(
			t,
			alpacaResponse.(map[string]interface{})["contact"].(map[string]interface{})["state"],
			"state",
		)

		assert.Equal(
			t,
			alpacaResponse.(map[string]interface{})["contact"].(map[string]interface{})["postal_code"],
			"postalCode",
		)

		w.WriteHeader(http.StatusCreated)
		_, errWrite := w.Write(jsonBytes)

		assert.Nil(t, errWrite)
	}))

	defer server.Close()

	alpacaClient := alpaca.NewBroker(
		"alpacaAPIKey",
		"alpacaAPISecret",
		server.URL,
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
}

func TestAlpacaListPositions(t *testing.T) {
	t.Parallel()

	accountID := "accountID"
	expectedPath := fmt.Sprint("/v1/trading/accounts/", accountID, "/positions")

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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, expectedPath)
		assert.Equal(t, r.Method, http.MethodGet)
		assert.NotNil(t, r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusOK)
		_, errWrite := w.Write(jsonBytes)

		assert.Nil(t, errWrite)
	}))

	defer server.Close()

	alpacaClient := alpaca.NewBroker(
		"alpacaAPIKey",
		"alpacaAPISecret",
		server.URL,
	)

	positions, errPositions := alpacaClient.ListPositions(accountID)

	assert.Nil(t, errPositions)
	assert.Equal(t, len(positions), 2)

	assert.Equal(t, positions[0].PositionID, "asset_id_1")
	assert.Equal(t, positions[0].Symbol, "symbol_1")
	assert.Equal(t, positions[0].MarketValue, 345.67)

	assert.Equal(t, positions[1].PositionID, "asset_id_2")
	assert.Equal(t, positions[1].Symbol, "symbol_2")
	assert.Equal(t, positions[1].MarketValue, 123.45)
}
