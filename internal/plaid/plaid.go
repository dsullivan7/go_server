package plaid

import (
	"context"
	"fmt"

	"github.com/plaid/plaid-go/plaid"
)

type IClient interface {
	CreateToken(userID string) (string, error)
	CreateProcessorToken(accessToken string, accountID string, accessor string) (string, error)
	GetAccessToken(publicToken string) (string, error)
	GetAccount(accessToken string) (string, string, error)
}

type Client struct {
	client      *plaid.APIClient
	redirectURI string
}

func NewClient(
	client *plaid.APIClient,
	redirectURI string,
) IClient {
	return &Client{
		client:      client,
		redirectURI: redirectURI,
	}
}

// CreateToken creates a plaid token to use in link.
func (pc *Client) CreateToken(userID string) (string, error) {
	ctx := context.Background()

	request := plaid.NewLinkTokenCreateRequest(
		"Sunburst",
		"en",
		[]plaid.CountryCode{plaid.COUNTRYCODE_US},
		*plaid.NewLinkTokenCreateRequestUser(userID),
	)
	request.SetRedirectUri(pc.redirectURI)
	request.SetProducts([]plaid.Products{plaid.PRODUCTS_AUTH})

	resp, _, err := pc.client.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()

	if err != nil {
		plaidErr, _ := plaid.ToPlaidError(err)

		return "", fmt.Errorf("plaid error %w: %s", err, plaidErr.ErrorMessage)
	}

	return resp.GetLinkToken(), nil
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
