package alpaca

import (
	"time"
	"errors"
	"bytes"
	"encoding/base64"
	"context"
	"encoding/json"
	"fmt"
	"go_server/internal/broker"
	goServerHTTP "go_server/internal/http"
	"net/http"
)

var AlpacaAPIError = errors.New("alpaca api error")

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
func (brkr *Broker) CreateAccount(
	givenName string,
	familyName string,
	dateOfBirth string,
	taxID string,
	emailAddress string,
	phoneNumber string,
	streetAddress string,
	city string,
	state string,
	postalCode string,
	fundingSource string,
	ipAddress string,
) (string, error) {
	context := context.Background()

	body := map[string]interface{}{
		"contact": map[string]interface{}{
			"email_address": emailAddress,
			"phone_number":  phoneNumber,
			"street_address":  []string{streetAddress},
			"city": city,
			"state": state,
			"postal_code": postalCode,
			"country": "USA",
		},
		"identity": map[string]interface{}{
			"given_name": givenName,
			"family_name": familyName,
			"date_of_birth": dateOfBirth,
			"country_of_tax_residency": "USA",
			"tax_id": taxID,
			"tax_id_type": "USA_SSN",
			"country_of_tax_residence": "USA",
			"funding_source": []string{fundingSource},
		},
		"disclosures":  map[string]interface{}{
			"is_control_person": false,
			"is_affiliated_exchange_or_finra": false,
			"is_politically_exposed": false,
			"immediate_family_exposed": false,
		},
		"agreements": []interface{}{
	    map[string]string{
	      "agreement": "margin_agreement",
	      "signed_at": time.Now().Format(time.RFC3339),
	      "ip_address": ipAddress,
	      "revision": "16.2021.05",
	    },
	    map[string]string{
	      "agreement": "account_agreement",
	      "signed_at": time.Now().Format(time.RFC3339),
	      "ip_address": ipAddress,
	      "revision": "16.2021.05",
	    },
	    map[string]string{
	      "agreement": "customer_agreement",
	      "signed_at": time.Now().Format(time.RFC3339),
	      "ip_address": ipAddress,
	      "revision": "16.2021.05",
	    },
	  },
	}

	jsonBytes, errMarshal := json.Marshal(body)

	if errMarshal != nil {
		return "", fmt.Errorf("failed to construct the request body: %w", errMarshal)
	}

	req, errReq := http.NewRequestWithContext(
		context,
		http.MethodPost,
		fmt.Sprint(brkr.alpacaAPIURL, "/v1/accounts"),
		bytes.NewReader(jsonBytes),
	)

	if errReq != nil {
		return "", fmt.Errorf("failed to create the request: %w", errReq)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(brkr.alpacaAPIKey, ":", brkr.alpacaAPISecret)))

	req.Header = http.Header{
		"Authorization": []string{fmt.Sprint("Basic ", authHeader)},
	}

	res, errRes := brkr.httpClient.Do(req)

	if errRes != nil {
		return "", fmt.Errorf("failed to get the response: %w", errRes)
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var alpacaResponse map[string]interface{}

	errDecode := decoder.Decode(&alpacaResponse)

	if errDecode != nil {
		return "", fmt.Errorf("failed to decode the response: %w", errDecode)
	}

	if (res.StatusCode != 200) {
		return "", fmt.Errorf("%w: %s", AlpacaAPIError, alpacaResponse["message"])
	}

	return alpacaResponse["id"].(string), nil
}

// DeleteAccount deactivates an active account.
func (brkr *Broker) DeleteAccount(accountID string) (error) {
	context := context.Background()

	req, errReq := http.NewRequestWithContext(
		context,
		http.MethodDelete,
		fmt.Sprint(brkr.alpacaAPIURL, "/v1/accounts/", accountID),
		nil,
	)

	if errReq != nil {
		return fmt.Errorf("failed to create the request: %w", errReq)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(brkr.alpacaAPIKey, ":", brkr.alpacaAPISecret)))

	req.Header = http.Header{
		"Authorization": []string{fmt.Sprint("Basic ", authHeader)},
	}

	res, errRes := brkr.httpClient.Do(req)

	if errRes != nil {
		return fmt.Errorf("failed to get the response: %w", errRes)
	}

	defer res.Body.Close()

	if (res.StatusCode != 204) {
		decoder := json.NewDecoder(res.Body)

		var alpacaResponse map[string]interface{}

		errDecode := decoder.Decode(&alpacaResponse)

		if errDecode != nil {
			return fmt.Errorf("failed to decode the response: %w", errDecode)
		}

		return fmt.Errorf("%w: %s", AlpacaAPIError, alpacaResponse["message"])
	}

	return nil
}

// CreateOrder creates an order for an account.
func (brkr *Broker) CreateOrder(accountID string, symbol string, quantity float32, side string) (string, error) {
	context := context.Background()

	body := map[string]interface{}{
		"symbol": symbol,
		"quantity":  quantity,
		"side":  side,
		"type": "market",
		"time_in_force": "day",
	}

	jsonBytes, errMarshal := json.Marshal(body)

	if errMarshal != nil {
		return "", fmt.Errorf("failed to construct the request body: %w", errMarshal)
	}

	req, errReq := http.NewRequestWithContext(
		context,
		http.MethodPost,
		fmt.Sprint(brkr.alpacaAPIURL, "/v1/trading/accounts/", accountID, "/orders"),
		bytes.NewReader(jsonBytes),
	)

	if errReq != nil {
		return "", fmt.Errorf("failed to create the request: %w", errReq)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(brkr.alpacaAPIKey, ":", brkr.alpacaAPISecret)))

	req.Header = http.Header{
		"Authorization": []string{fmt.Sprint("Basic ", authHeader)},
	}

	res, errRes := brkr.httpClient.Do(req)

	if errRes != nil {
		return "", fmt.Errorf("failed to get the response: %w", errRes)
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var alpacaResponse map[string]interface{}

	errDecode := decoder.Decode(&alpacaResponse)

	if errDecode != nil {
		return "", fmt.Errorf("failed to decode the response: %w", errDecode)
	}

	if (res.StatusCode != 200) {
		return "", fmt.Errorf("%w: %s", AlpacaAPIError, alpacaResponse["message"])
	}

	return alpacaResponse["id"].(string), nil
}
