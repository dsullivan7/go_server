package plaid

import (
  "fmt"
  "context"

  "github.com/plaid/plaid-go/plaid"
)

type Client interface {
	CreatePlaidToken(userID string) (string, error)
	ExchangePublicToken(publicToken string) (string, string, string, error)
}

type Implementation struct {
	plaidClient *plaid.APIClient
}

func NewPlaidClient(
	plaidClient *plaid.APIClient,
) Client {
	return &Implementation{
		plaidClient: plaidClient,
	}
}

func (pc *Implementation) CreatePlaidToken(userID string) (string, error) {
  ctx := context.Background()

  request := plaid.NewLinkTokenCreateRequest("Sunburst", "en", []plaid.CountryCode{plaid.COUNTRYCODE_US}, *plaid.NewLinkTokenCreateRequestUser(userID))
  request.SetRedirectUri("http://localhost:3000")
  request.SetProducts([]plaid.Products{plaid.PRODUCTS_AUTH})

  resp, _, err := pc.plaidClient.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()

  if (err != nil) {
    plaidErr, _ := plaid.ToPlaidError(err)

    return "", fmt.Errorf("plaid error %w: %s", err, plaidErr.ErrorMessage)
  }

  return resp.GetLinkToken(), nil
}

func (pc *Implementation) ExchangePublicToken(publicToken string) (string, string, string, error) {
  ctx := context.Background()

  exchangePublicTokenReq := plaid.NewItemPublicTokenExchangeRequest(publicToken)
  exchangePublicTokenResp, _, errAccessToken := pc.plaidClient.PlaidApi.ItemPublicTokenExchange(ctx).
    ItemPublicTokenExchangeRequest(*exchangePublicTokenReq).
    Execute()

  if (errAccessToken != nil) {
    plaidErr, _ := plaid.ToPlaidError(errAccessToken)

    return "", "", "", fmt.Errorf("plaid error %w: %s", errAccessToken, plaidErr.ErrorMessage)
  }

  accessToken := exchangePublicTokenResp.GetAccessToken()
  itemID := exchangePublicTokenResp.GetItemId()

  accountsGetResp, _, errAccount := pc.plaidClient.PlaidApi.AccountsGet(ctx).
    AccountsGetRequest(*plaid.NewAccountsGetRequest(accessToken)).
    Execute()

  if (errAccount != nil) {
    plaidErr, _ := plaid.ToPlaidError(errAccount)

    return "", "", "", fmt.Errorf("plaid error %w: %s", errAccount, plaidErr.ErrorMessage)
  }

  institutionGetResp, _, errInstitution := pc.plaidClient.PlaidApi.InstitutionsGetById(ctx).
    InstitutionsGetByIdRequest(
      *plaid.NewInstitutionsGetByIdRequest(*accountsGetResp.Item.InstitutionId.Get(),
      []plaid.CountryCode{plaid.COUNTRYCODE_US},
    )).
    Execute()

  if (errInstitution != nil) {
    plaidErr, _ := plaid.ToPlaidError(errInstitution)

    return "", "", "", fmt.Errorf("plaid error %w: %s", errInstitution, plaidErr.ErrorMessage)
  }

  return accessToken, itemID, institutionGetResp.Institution.Name, nil
}
