package dwolla

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go_server/internal/bank"
	"go_server/internal/logger"
	"go_server/internal/models"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ErrDwollaAPI = errors.New("dwolla api error")
var ErrTypeConversion = errors.New("type conversion error")

const centsToDollars = 100

type Bank struct {
	dwollaAPIAccessToken          string
	dwollaAPIAccessTokenExpiresAt time.Time
	dwollaAPIKey                  string
	dwollaAPISecret               string
	dwollaAPIURL                  string
	dwollaWebhookURL              string
	dwollaWebhookSecret           string
	logger                        logger.Logger
}

func NewBank(
	dwollaAPIKey string,
	dwollaAPISecret string,
	dwollaAPIURL string,
	dwollaWebhookURL string,
	dwollaWebhookSecret string,
	lggr logger.Logger,
) bank.Bank {
	return &Bank{
		dwollaAPIKey:        dwollaAPIKey,
		dwollaAPISecret:     dwollaAPISecret,
		dwollaAPIURL:        dwollaAPIURL,
		dwollaWebhookURL:    dwollaWebhookURL,
		dwollaWebhookSecret: dwollaWebhookSecret,
		logger:              lggr,
	}
}

//nolint:funlen,cyclop,unparam
func (bnk *Bank) sendRequest(
	path string,
	method string,
	body map[string]interface{},
) (interface{}, error) {
	if bnk.dwollaAPIAccessTokenExpiresAt.IsZero() || time.Now().After(bnk.dwollaAPIAccessTokenExpiresAt) {
		errAuth := bnk.authenticate()
		if errAuth != nil {
			return nil, fmt.Errorf("failed to authenticate: %w", errAuth)
		}
	}

	ctxt := context.Background()

	jsonBytes, errMarshal := json.Marshal(body)

	if errMarshal != nil {
		return nil, fmt.Errorf("failed to construct the request body: %w", errMarshal)
	}

	req, errReq := http.NewRequestWithContext(
		ctxt,
		method,
		fmt.Sprint(bnk.dwollaAPIURL, path),
		bytes.NewReader(jsonBytes),
	)

	if errReq != nil {
		return nil, fmt.Errorf("failed to create the request: %w", errReq)
	}

	req.Header.Set("Accept", "application/vnd.dwolla.v1.hal+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", bnk.dwollaAPIAccessToken))

	res, errRes := http.DefaultClient.Do(req)

	if errRes != nil {
		return nil, fmt.Errorf("failed to get the response: %w", errRes)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusCreated {
		return res.Header.Get("Location"), nil
	}

	decoder := json.NewDecoder(res.Body)

	var dwollaResponse interface{}

	errDecode := decoder.Decode(&dwollaResponse)

	if errDecode != nil {
		return nil, fmt.Errorf("failed to decode the response: %w", errDecode)
	}

	if res.StatusCode != http.StatusOK &&
		res.StatusCode != http.StatusCreated &&
		res.StatusCode != http.StatusNoContent {
		bnk.logger.InfoWithMeta("Dwolla debug", map[string]interface{}{"dwollaResponse": dwollaResponse})

		return nil, fmt.Errorf("%w: %s", ErrDwollaAPI, dwollaResponse.(map[string]interface{})["message"].(string))
	}

	return dwollaResponse, nil
}

func (bnk *Bank) authenticate() error {
	ctxt := context.Background()

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, errReq := http.NewRequestWithContext(
		ctxt,
		http.MethodPost,
		fmt.Sprint(bnk.dwollaAPIURL, "/token"),
		strings.NewReader(data.Encode()),
	)

	if errReq != nil {
		return fmt.Errorf("failed to create the request: %w", errReq)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(bnk.dwollaAPIKey, ":", bnk.dwollaAPISecret)))

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprint("Basic ", authHeader))

	res, errRes := http.DefaultClient.Do(req)

	if errRes != nil {
		return fmt.Errorf("failed to get the response: %w", errRes)
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var dwollaResponse interface{}

	errDecode := decoder.Decode(&dwollaResponse)

	if errDecode != nil {
		return fmt.Errorf("failed to decode the response: %w", errDecode)
	}

	if res.StatusCode != http.StatusOK &&
		res.StatusCode != http.StatusCreated &&
		res.StatusCode != http.StatusNoContent {
		bnk.logger.InfoWithMeta("Dwolla auth debug", map[string]interface{}{"dwollaResponse": dwollaResponse})

		return fmt.Errorf("%w: %s", ErrDwollaAPI, dwollaResponse.(map[string]interface{})["message"].(string))
	}

	expirationDuration := time.Second * time.Duration(int(dwollaResponse.(map[string]interface{})["expires_in"].(float64)))

	apiAccessToken, ok := dwollaResponse.(map[string]interface{})["access_token"].(string)

	if !ok {
		return ErrTypeConversion
	}

	bnk.dwollaAPIAccessToken = apiAccessToken
	bnk.dwollaAPIAccessTokenExpiresAt = time.Now().Add(expirationDuration)

	return nil
}

// CreateTransfer creates a transfer for the source to the destination.
func (bnk *Bank) CreateTransfer(
	source models.BankAccount,
	destination models.BankAccount,
	amount int,
) (*models.BankTransfer, error) {
	body := map[string]interface{}{
		"_links": map[string]interface{}{
			"source": map[string]string{
				"href": fmt.Sprint(bnk.dwollaAPIURL, "/funding-sources/", *source.DwollaFundingSourceID),
			},
			"destination": map[string]string{
				"href": fmt.Sprint(bnk.dwollaAPIURL, "/funding-sources/", *destination.DwollaFundingSourceID),
			},
		},
		"amount": map[string]string{
			"value":    fmt.Sprintf("%f", float64(amount)/centsToDollars),
			"currency": "USD",
		},
	}

	dwollaResponse, errDwolla := bnk.sendRequest(
		"/transfers",
		http.MethodPost,
		body,
	)

	if errDwolla != nil {
		return nil, errDwolla
	}

	split := strings.Split(dwollaResponse.(string), "/")
	dwollaTransferID := split[len(split)-1]

	return &models.BankTransfer{DwollaTransferID: &dwollaTransferID}, nil
}

// CreateCustomer creates a customer within dwolla.
func (bnk *Bank) CreateCustomer(user models.User) (*models.User, error) {
	if user.DwollaCustomerID != nil {
		// this user already has a customer, return the user
		return &user, nil
	}

	body := map[string]interface{}{
		"firstName":   user.FirstName,
		"lastName":    user.LastName,
		"email":       user.Email,
		"type":        "personal",
		"address1":    user.Address,
		"city":        user.City,
		"state":       user.State,
		"postalCode":  user.PostalCode,
		"dateOfBirth": user.DateOfBirth,
		"ssn":         user.SSN,
	}

	dwollaResponse, errDwolla := bnk.sendRequest(
		"/customers",
		http.MethodPost,
		body,
	)

	if errDwolla != nil {
		return nil, errDwolla
	}

	split := strings.Split(dwollaResponse.(string), "/")
	dwollaCustomerID := split[len(split)-1]

	return &models.User{DwollaCustomerID: &dwollaCustomerID}, nil
}

// CreateBank creates a funding source within dwolla.
func (bnk *Bank) CreateBankAccount(user models.User, plaidProcessorToken string) (*models.BankAccount, error) {
	body := map[string]interface{}{
		"name":       user.UserID.String(),
		"plaidToken": plaidProcessorToken,
	}

	dwollaResponse, errDwolla := bnk.sendRequest(
		fmt.Sprint("/customers/", *user.DwollaCustomerID, "/funding-sources"),
		http.MethodPost,
		body,
	)

	if errDwolla != nil {
		return nil, errDwolla
	}

	split := strings.Split(dwollaResponse.(string), "/")
	dwollaFundingSourceID := split[len(split)-1]

	return &models.BankAccount{DwollaFundingSourceID: &dwollaFundingSourceID}, nil
}

// CreateWebhook creates a webhook for dwolla.
func (bnk *Bank) CreateWebhook() (*models.Webhook, error) {
	body := map[string]interface{}{
		"url":    bnk.dwollaWebhookURL,
		"secret": bnk.dwollaWebhookSecret,
	}

	dwollaResponse, errDwolla := bnk.sendRequest(
		"/webhook-subscriptions",
		http.MethodPost,
		body,
	)

	if errDwolla != nil {
		return nil, errDwolla
	}

	split := strings.Split(dwollaResponse.(string), "/")
	dwollaID := split[len(split)-1]

	return &models.Webhook{DwollaWebhookID: &dwollaID}, nil
}

// GetPlaidAccessor returns the accessor for plaid access tokens.
func (bnk *Bank) GetPlaidAccessor() string {
	return "dwolla"
}
