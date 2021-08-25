package errors

import (
	"net/http"
)

type HTTPServerError struct {
	HTTPStatus int
	Message    string
	Err        error
}

func (err HTTPServerError) GetHTTPStatus() int {
	if err.HTTPStatus == 0 {
		return http.StatusInternalServerError
	}

	return err.HTTPStatus
}

func (err HTTPServerError) GetMessage() string {
	if err.Message == "" {
		return "Something went wrong, we're working to resolve the problem."
	}

	return err.Message
}

func (err HTTPServerError) GetError() error {
	return err.Err
}
