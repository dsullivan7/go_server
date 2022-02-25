package alpaca

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go_server/internal/broker"
	goServerHTTP "go_server/internal/http"
	"net/http"
	"strconv"
	"time"
)

var ErrAlpacaAPI = errors.New("alpaca api error")

const bitConversion = 64

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

func (brkr *Broker) sendRequest(
	path string,
	method string,
	body map[string]interface{},
) (interface{}, error) {
	context := context.Background()

	jsonBytes, errMarshal := json.Marshal(body)

	if errMarshal != nil {
		return nil, fmt.Errorf("failed to construct the request body: %w", errMarshal)
	}

	req, errReq := http.NewRequestWithContext(
		context,
		method,
		fmt.Sprint(brkr.alpacaAPIURL, path),
		bytes.NewReader(jsonBytes),
	)

	if errReq != nil {
		return nil, fmt.Errorf("failed to create the request: %w", errReq)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(brkr.alpacaAPIKey, ":", brkr.alpacaAPISecret)))

	req.Header = http.Header{
		"Authorization": []string{fmt.Sprint("Basic ", authHeader)},
	}

	res, errRes := brkr.httpClient.Do(req)

	if errRes != nil {
		return nil, fmt.Errorf("failed to get the response: %w", errRes)
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var alpacaResponse interface{}

	errDecode := decoder.Decode(&alpacaResponse)

	if errDecode != nil {
		return nil, fmt.Errorf("failed to decode the response: %w", errDecode)
	}

	if res.StatusCode != http.StatusOK &&
		res.StatusCode != http.StatusCreated &&
		res.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("%w: %s", ErrAlpacaAPI, alpacaResponse.(map[string]interface{})["message"].(string))
	}

	return alpacaResponse, nil
}

func getCreateAccountBody(
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
) map[string]interface{} {
	return map[string]interface{}{
		"contact": map[string]interface{}{
			"email_address":  emailAddress,
			"phone_number":   phoneNumber,
			"street_address": []string{streetAddress},
			"city":           city,
			"state":          state,
			"postal_code":    postalCode,
			"country":        "USA",
		},
		"identity": map[string]interface{}{
			"given_name":               givenName,
			"family_name":              familyName,
			"date_of_birth":            dateOfBirth,
			"country_of_tax_residency": "USA",
			"tax_id":                   taxID,
			"tax_id_type":              "USA_SSN",
			"country_of_tax_residence": "USA",
			"funding_source":           []string{fundingSource},
		},
		"disclosures": map[string]interface{}{
			"is_control_person":               false,
			"is_affiliated_exchange_or_finra": false,
			"is_politically_exposed":          false,
			"immediate_family_exposed":        false,
		},
		"agreements": []interface{}{
			map[string]string{
				"agreement":  "margin_agreement",
				"signed_at":  time.Now().Format(time.RFC3339),
				"ip_address": ipAddress,
				"revision":   "16.2021.05",
			},
			map[string]string{
				"agreement":  "account_agreement",
				"signed_at":  time.Now().Format(time.RFC3339),
				"ip_address": ipAddress,
				"revision":   "16.2021.05",
			},
			map[string]string{
				"agreement":  "customer_agreement",
				"signed_at":  time.Now().Format(time.RFC3339),
				"ip_address": ipAddress,
				"revision":   "16.2021.05",
			},
		},
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
	body := getCreateAccountBody(
		givenName,
		familyName,
		dateOfBirth,
		taxID,
		emailAddress,
		phoneNumber,
		streetAddress,
		city,
		state,
		postalCode,
		fundingSource,
		ipAddress,
	)

	alpacaResponse, errAlpaca := brkr.sendRequest(
		"/v1/accounts",
		http.MethodPost,
		body,
	)

	if errAlpaca != nil {
		return "", errAlpaca
	}

	return alpacaResponse.(map[string]interface{})["id"].(string), nil
}

// GetAccount retreives the given account.
func (brkr *Broker) GetAccount(accountID string) (*broker.Account, error) {
	alpacaResponse, errAlpaca := brkr.sendRequest(
		fmt.Sprint("/v1/trading/accounts/", accountID, "/account"),
		http.MethodGet,
		nil,
	)

	if errAlpaca != nil {
		return nil, errAlpaca
	}

	cash, errCash := strconv.ParseFloat(alpacaResponse.(map[string]interface{})["cash"].(string), bitConversion)

	if errCash != nil {
		return nil, fmt.Errorf("error parsing the cash amount: %w", errCash)
	}

	account := broker.Account{
		AccountID: alpacaResponse.(map[string]interface{})["id"].(string),
		Cash:      cash,
	}

	return &account, nil
}

// ListAccounts retreives a list of accounts given the query.
func (brkr *Broker) ListAccounts(query string) ([]broker.Account, error) {
	alpacaResponse, errAlpaca := brkr.sendRequest(
		fmt.Sprint("/v1/accounts?query=", query),
		http.MethodGet,
		nil,
	)

	if errAlpaca != nil {
		return nil, errAlpaca
	}

	accounts := make([]broker.Account, len(alpacaResponse.([]interface{})))

	for i, alpacaAccount := range alpacaResponse.([]interface{}) {
		cash, errCash := strconv.ParseFloat(alpacaAccount.(map[string]interface{})["cash"].(string), bitConversion)

		if errCash != nil {
			return nil, fmt.Errorf("error parsing the cash amount: %w", errCash)
		}

		accounts[i] = broker.Account{
			AccountID: alpacaAccount.(map[string]interface{})["id"].(string),
			Cash:      cash,
		}
	}

	return accounts, nil
}

// DeleteAccount deactivates an active account.
func (brkr *Broker) DeleteAccount(accountID string) error {
	_, errAlpaca := brkr.sendRequest(
		fmt.Sprint("/v1/accounts/", accountID),
		http.MethodDelete,
		nil,
	)

	if errAlpaca != nil {
		return errAlpaca
	}

	return nil
}

// CreateOrder creates an order for an account.
func (brkr *Broker) CreateOrder(accountID string, symbol string, amount float64, side string) (string, error) {
	body := map[string]interface{}{
		"symbol":        symbol,
		"notional": 		 amount,
		"side":          side,
		"type":          "market",
		"time_in_force": "day",
	}

	alpacaResponse, errAlpaca := brkr.sendRequest(
		fmt.Sprint("/v1/trading/accounts/", accountID, "/orders"),
		http.MethodPost,
		body,
	)

	if errAlpaca != nil {
		return "", errAlpaca
	}

	return alpacaResponse.(map[string]interface{})["id"].(string), nil
}

// CreateTransfer creates a transfer for an account.
func (brkr *Broker) CreateTransfer(
	accountID string,
	relationshipID string,
	amount float64,
	direction string,
) (string, error) {
	body := map[string]interface{}{
		"transfer_type":   "ach",
		"relationship_id": relationshipID,
		"amount":          amount,
		"direction":       direction,
	}

	alpacaResponse, errAlpaca := brkr.sendRequest(
		fmt.Sprint("/v1/accounts/", accountID, "/transfers"),
		http.MethodPost,
		body,
	)

	if errAlpaca != nil {
		return "", errAlpaca
	}

	return alpacaResponse.(map[string]interface{})["id"].(string), nil
}

// CreateACHRelationship creates an ach relastionship for an account.
func (brkr *Broker) CreateACHRelationship(
	accountID string,
	processorToken string,
) (string, error) {
	body := map[string]interface{}{
		"processor_token": processorToken,
	}

	alpacaResponse, errAlpaca := brkr.sendRequest(
		fmt.Sprint("/v1/accounts/", accountID, "/ach_relationships"),
		http.MethodPost,
		body,
	)

	if errAlpaca != nil {
		return "", errAlpaca
	}

	return alpacaResponse.(map[string]interface{})["id"].(string), nil
}
