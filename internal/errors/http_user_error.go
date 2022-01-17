package errors

import (
	"net/http"
)

type HTTPUserError struct {
	HTTPStatus int
	Message    string
	Code       string
	Err        error
}

func (err HTTPUserError) GetHTTPStatus() int {
	if err.HTTPStatus == 0 {
		return http.StatusBadRequest
	}

	return err.HTTPStatus
}

func (err HTTPUserError) GetMessage() string {
	if err.Message == "" {
		return "Something went wrong, try again."
	}

	return err.Message
}

func (err HTTPUserError) GetCode() string {
	if err.Code == "" {
		return "bad_request"
	}

	return err.Code
}

func (err HTTPUserError) GetError() error {
	return err.Err
}
