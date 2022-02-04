package http

import (
	"fmt"
	"net/http"
)

type IClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
}

func NewClient() IClient {
	return &Client{}
}

// Do completes the http request.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error doing the request: %w", err)
	}

	return res, nil
}
