package http

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func NewClient() *MockClient {
	return &MockClient{}
}

func (mockClient *MockClient) Do(req *http.Request) (*http.Response, error) {
	args := mockClient.Called(req)

	return args.Get(0).(*http.Response), args.Error(1)
}
