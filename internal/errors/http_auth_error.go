package errors

import (
	"net/http"
)

type HTTPAuthError struct {
	HTTPStatus int
	Message    string
	Err        error
}

func (err HTTPAuthError) GetHTTPStatus() int {
	if err.HTTPStatus == 0 {
		return http.StatusForbidden
	}

	return err.HTTPStatus
}

func (err HTTPAuthError) GetMessage() string {
	if err.Message == "" {
		return "There was an authorization error."
	}

	return err.Message
}

func (err HTTPAuthError) GetError() error {
	return err.Err
}
