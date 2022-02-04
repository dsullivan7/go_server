package alpaca

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go_server/internal/broker"
	goServerHTTP "go_server/internal/http"
	"net/http"
)

type Account struct {
}

type Broker struct {
	alpacaAPIKey    string
	alpacaAPISecret string
	alpacaAPIURL    string
	httpClient      goServerHTTP.IClient
}

func NewBroker(
	alpacaAPIKey string,
	alpacaAPISecret string,
	alpacaAPIURL string,
	httpClient goServerHTTP.IClient,
) broker.Broker {
	return &Broker{
		alpacaAPIKey:    alpacaAPIKey,
		alpacaAPISecret: alpacaAPISecret,
		alpacaAPIURL:    alpacaAPIURL,
		httpClient:      httpClient,
	}
}

// CreateAccount creates an account for the user.
func (b *Broker) CreateAccount(
	emailAddress string,
	phoneNumber string,
) (string, error) {
	context := context.Background()

	body := map[string]interface{}{
		"contact": map[string]interface{}{
			"email_address": emailAddress,
			"phone_number":  phoneNumber,
		},
	}

	jsonBytes, errMarshal := json.Marshal(body)

	if errMarshal != nil {
		return "", fmt.Errorf("failed to construct the request body: %w", errMarshal)
	}

	req, errReq := http.NewRequestWithContext(
		context,
		http.MethodPost,
		fmt.Sprint(b.alpacaAPIURL, "/v1/accounts"),
		bytes.NewReader(jsonBytes),
	)

	if errReq != nil {
		return "", fmt.Errorf("failed to create the request: %w", errReq)
	}

	res, errRes := b.httpClient.Do(req)

	if errRes != nil {
		return "", fmt.Errorf("failed to get the response: %w", errRes)
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var alpacaResponse map[string]map[string]string

	errDecode := decoder.Decode(&alpacaResponse)

	if errDecode != nil {
		return "", fmt.Errorf("failed to decode the response: %w", errDecode)
	}

	return alpacaResponse["account"]["id"], nil
}
