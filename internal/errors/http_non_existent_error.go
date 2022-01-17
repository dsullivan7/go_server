package errors

import (
	"errors"
	"net/http"
)

type HTTPNonExistentError struct{}

const message = "this record does not exist"

var errNonExistent = errors.New(message)

func (err HTTPNonExistentError) GetHTTPStatus() int {
	return http.StatusBadRequest
}

func (err HTTPNonExistentError) GetMessage() string {
	return "This record does not exist."
}

func (err HTTPNonExistentError) GetCode() string {
	return "non_existent"
}

func (err HTTPNonExistentError) GetError() error {
	return errNonExistent
}
