package plaid

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go_server/internal/logger"

	"github.com/plaid/plaid-go/plaid"
)

var ErrPlaidAPI = errors.New("plaid api error")
var ErrTypeConversionError = errors.New("type conversion error")

type IClient interface {
	CreateToken(userID string) (string, error)
	CreateProcessorToken(accessToken string, accountID string, accessor string) (string, error)
	GetAccessToken(publicToken string) (string, error)
	GetAccount(accessToken string) (string, string, error)
	CreateTransferAuthorization(
		accessToken string,
		accountID string,
		originationAccountID string,
		amount string,
		transferType string,
		legalName string,
	) (string, error)
	CreateTransfer(
		accessToken string,
		accountID string,
		originationAccountID string,
		authorizationID string,
		amount string,
		transferType string,
		legalName string,
	) (string, error)
}

type Client struct {
	clientID    string
	secret      string
	apiURL      string
	redirectURI string
	logger      logger.Logger
}

func NewClient(
	clientID string,
	secret string,
	apiURL string,
	redirectURI string,
	lggr logger.Logger,
) IClient {
	return &Client{
		clientID:    clientID,
		secret:      secret,
		apiURL:      apiURL,
		redirectURI: redirectURI,
		logger:      lggr,
	}
}

func (pc *Client) sendRequest(
	path string,
	body map[string]interface{},
) (interface{}, error) {
	method := http.MethodPost
	context := context.Background()

	jsonBytes, errMarshal := json.Marshal(body)

	if errMarshal != nil {
		return nil, fmt.Errorf("failed to construct the request body: %w", errMarshal)
	}

	req, errReq := http.NewRequestWithContext(
		context,
		method,
		fmt.Sprint(pc.apiURL, path),
		bytes.NewReader(jsonBytes),
	)

	if errReq != nil {
		return nil, fmt.Errorf("failed to create the request: %w", errReq)
	}

	req.Header.Set("PLAID-CLIENT-ID", pc.clientID)
	req.Header.Set("PLAID-SECRET", pc.secret)
	req.Header.Set("Content-Type", "application/json")

	res, errRes := http.DefaultClient.Do(req)

	if errRes != nil {
		return nil, fmt.Errorf("failed to get the response: %w", errRes)
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var plaidResponse interface{}

	errDecode := decoder.Decode(&plaidResponse)

	if errDecode != nil {
		return nil, fmt.Errorf("failed to decode the response: %w", errDecode)
	}

	if res.StatusCode != http.StatusOK &&
		res.StatusCode != http.StatusCreated &&
		res.StatusCode != http.StatusNoContent {
		pc.logger.InfoWithMeta("Plaid debug", map[string]interface{}{"plaidResponse": plaidResponse})

		return nil, fmt.Errorf("%w: %s", ErrPlaidAPI, plaidResponse.(map[string]interface{})["message"].(string))
	}

	return plaidResponse, nil
}

// CreateToken creates a plaid token to use in link.
func (pc *Client) CreateToken(userID string) (string, error) {
	tokenResp, tokenError := pc.sendRequest(
		"/link/token/create",
		map[string]interface{}{
			"user": map[string]string{
				"client_user_id": userID,
			},
			"client_name":   "Sunburst",
			"products":      []string{string(plaid.PRODUCTS_AUTH)},
			"country_codes": []string{string(plaid.COUNTRYCODE_US)},
			"language":      "en",
			"redirect_uri":  pc.redirectURI,
			"account_filters": map[string]interface{}{
				"depository": map[string]interface{}{
					"account_subtypes": []string{"checking"},
				},
			},
		},
	)

	if tokenError != nil {
		return "", tokenError
	}

	return tokenResp.(map[string]interface{})["link_token"].(string), nil
}

// CreateTransferAuthorization creates an authorization for a bank transfer.
func (pc *Client) CreateTransferAuthorization(
	accessToken string,
	accountID string,
	originationAccountID string,
	amount string,
	transferType string,
	legalName string,
) (string, error) {
	transferAuthResp, transferAuthError := pc.sendRequest(
		"/transfer/authorization/create",
		map[string]interface{}{
			"access_token":           accessToken,
			"account_id":             accountID,
			"origination_account_id": originationAccountID,
			"amount":                 amount,
			"type":                   transferType,
			"network":                "ach",
			"ach_class":              "web",
			"user": map[string]interface{}{
				"legal_name": legalName,
			},
		},
	)

	if transferAuthError != nil {
		return "", transferAuthError
	}

	return transferAuthResp.(map[string]interface{})["authorization"].(map[string]interface{})["id"].(string), nil
}

// CreateTransfer creates a bank transfer.
func (pc *Client) CreateTransfer(
	accessToken string,
	accountID string,
	originationAccountID string,
	authorizationID string,
	amount string,
	transferType string,
	legalName string,
) (string, error) {
	transferAuthResp, transferAuthError := pc.sendRequest(
		"/transfer/create",
		map[string]interface{}{
			"authorization_id":       authorizationID,
			"access_token":           accessToken,
			"account_id":             accountID,
			"origination_account_id": originationAccountID,
			"amount":                 amount,
			"type":                   transferType,
			"network":                "ach",
			"ach_class":              "web",
			"user": map[string]interface{}{
				"legal_name": legalName,
			},
		},
	)

	if transferAuthError != nil {
		return "", transferAuthError
	}

	return transferAuthResp.(map[string]interface{})["transfer"].(map[string]interface{})["id"].(string), nil
}

// GetAccessToken exchanges the plaid token for an access token.
func (pc *Client) GetAccessToken(publicToken string) (string, error) {
	exchangePublicTokenResp, errAccessToken := pc.sendRequest(
		"/item/public_token/exchange",
		map[string]interface{}{
			"public_token": publicToken,
		},
	)

	if errAccessToken != nil {
		return "", errAccessToken
	}

	return exchangePublicTokenResp.(map[string]interface{})["access_token"].(string), nil
}

// GetAccount uses the access token to retrieve the account information.
func (pc *Client) GetAccount(accessToken string) (string, string, error) {
	accountsGetResp, errAccount := pc.sendRequest(
		"/accounts/get",
		map[string]interface{}{
			"access_token": accessToken,
		},
	)

	if errAccount != nil {
		return "", "", errAccount
	}

	// nolint: lll
	accountID, okAccountType := accountsGetResp.(map[string]interface{})["accounts"].([]interface{})[0].(map[string]interface{})["account_id"].(string)

	if !okAccountType {
		return "", "", fmt.Errorf("accountID: %w", ErrTypeConversionError)
	}

	// nolint: lll
	institutionID, okInstType := accountsGetResp.(map[string]interface{})["item"].(map[string]interface{})["institution_id"].(string)

	if !okInstType {
		return "", "", fmt.Errorf("institutionID: %w", ErrTypeConversionError)
	}

	institutionGetResp, errInstitution := pc.sendRequest(
		"/institutions/get_by_id",
		map[string]interface{}{
			"institution_id": institutionID,
			"country_codes":  []string{string(plaid.COUNTRYCODE_US)},
		},
	)

	if errInstitution != nil {
		return "", "", errInstitution
	}

	// nolint: lll
	institutionName, okInstNameType := institutionGetResp.(map[string]interface{})["institution"].(map[string]interface{})["name"].(string)

	if !okInstNameType {
		return "", "", fmt.Errorf("institutionName: %w", ErrTypeConversionError)
	}

	return accountID, institutionName, nil
}

// CreateToken creates a processor token to use in link.
func (pc *Client) CreateProcessorToken(accessToken string, accountID string, processor string) (string, error) {
	processorTokenResp, errProcessorToken := pc.sendRequest(
		"/processor/token/create",
		map[string]interface{}{
			"access_token": accessToken,
			"account_id":   accountID,
			"processor":    processor,
		},
	)

	if errProcessorToken != nil {
		return "", errProcessorToken
	}

	return processorTokenResp.(map[string]interface{})["processor_token"].(string), nil
}
