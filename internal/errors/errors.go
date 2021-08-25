package errors

import (
	"net/http"
)

type HTTPError interface {
	GetHTTPStatus() int
	GetMessage() string
	GetError() error
}

type HTTPUserError struct {
	HTTPStatus int    `defaults:"400"`
	Message    string `defaults:"Something went wrong, try again."`
	Err      error
}

func (err HTTPUserError) GetHTTPStatus() int {
	if (err.HTTPStatus == 0) {
		return http.StatusBadRequest
	}

	return err.HTTPStatus
}

func (err HTTPUserError) GetMessage() string {
	if (err.Message == "") {
		return "an error occurred"
	}

	return err.Message
}

func (err HTTPUserError) GetError() error {
	return err.Err
}
