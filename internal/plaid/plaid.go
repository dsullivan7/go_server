package plaid

import (
  "fmt"
  "context"

  "github.com/plaid/plaid-go/plaid"
)

type Client interface {
	CreatePlaidToken(userID string) (string, error)
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
