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
	client      *plaid.APIClient
	apiURL      string
	redirectURI string
	clientID    string
	secret      string
	logger      logger.Logger
}

func NewClient(
	client *plaid.APIClient,
	clientID string,
	secret string,
	apiURL string,
	redirectURI string,
	lggr logger.Logger,
) IClient {
	return &Client{
		client:      client,
		redirectURI: redirectURI,
		clientID:    clientID,
		secret:      secret,
		apiURL:      apiURL,
		logger:      lggr,
	}
}

func (pc *Client) sendRequest(
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
		http.MethodPost,
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
		http.MethodPost,
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
		http.MethodPost,
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
	ctx := context.Background()

	exchangePublicTokenReq := plaid.NewItemPublicTokenExchangeRequest(publicToken)
	exchangePublicTokenResp, _, errAccessToken := pc.client.PlaidApi.ItemPublicTokenExchange(ctx).
		ItemPublicTokenExchangeRequest(*exchangePublicTokenReq).
		Execute()

	if errAccessToken != nil {
		plaidErr, _ := plaid.ToPlaidError(errAccessToken)

		return "", fmt.Errorf("plaid error %w: %s", errAccessToken, plaidErr.ErrorMessage)
	}

	accessToken := exchangePublicTokenResp.GetAccessToken()

	return accessToken, nil
}

// GetAccount uses the access token to retrieve the account information.
func (pc *Client) GetAccount(accessToken string) (string, string, error) {
	ctx := context.Background()

	accountsGetResp, _, errAccount := pc.client.PlaidApi.AccountsGet(ctx).
		AccountsGetRequest(*plaid.NewAccountsGetRequest(accessToken)).
		Execute()

	if errAccount != nil {
		plaidErr, _ := plaid.ToPlaidError(errAccount)

		return "", "", fmt.Errorf("plaid error %w: %s", errAccount, plaidErr.ErrorMessage)
	}

	institutionGetResp, _, errInstitution := pc.client.PlaidApi.InstitutionsGetById(ctx).
		InstitutionsGetByIdRequest(
			*plaid.NewInstitutionsGetByIdRequest(*accountsGetResp.Item.InstitutionId.Get(),
				[]plaid.CountryCode{plaid.COUNTRYCODE_US},
			)).
		Execute()

	if errInstitution != nil {
		plaidErr, _ := plaid.ToPlaidError(errInstitution)

		return "", "", fmt.Errorf("plaid error %w: %s", errInstitution, plaidErr.ErrorMessage)
	}

	institution := institutionGetResp.GetInstitution()

	return accountsGetResp.GetAccounts()[0].GetAccountId(), institution.GetName(), nil
}

// CreateToken creates a plaid token to use in link.
func (pc *Client) CreateProcessorToken(accessToken string, accountID string, processor string) (string, error) {
	ctx := context.Background()

	request := plaid.NewProcessorTokenCreateRequest(
		accessToken,
		accountID,
		processor,
	)

	resp, _, err := pc.client.PlaidApi.ProcessorTokenCreate(ctx).ProcessorTokenCreateRequest(*request).Execute()

	if err != nil {
		plaidErr, _ := plaid.ToPlaidError(err)

		return "", fmt.Errorf("plaid error %w: %s", err, plaidErr.ErrorMessage)
	}

	return resp.GetProcessorToken(), nil
}
