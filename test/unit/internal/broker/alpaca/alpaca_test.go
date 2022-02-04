package zap_test

import (
	"bytes"
	"go_server/internal/broker/alpaca"
	mockHTTP "go_server/test/mocks/http"
  "encoding/json"
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
		"account": map[string]string{
      "id": "test",
    },
	}

	jsonBytes, errMarshal := json.Marshal(body)

  assert.Nil(t, errMarshal)

	mockHTTPClient.On("Do", mock.Anything).Return(
    &http.Response{
      StatusCode: 200,
      Body: ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))),
    },
    nil,
  )

	accountNumber, errAcc := alpacaClient.CreateAccount(
		"emailAddress",
		"phoneNumber",
	)

	assert.Nil(t, errAcc)
	assert.Equal(t, accountNumber, "test")

	mockHTTPClient.AssertExpectations(t)
}
