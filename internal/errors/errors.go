package errors

type HTTPError interface {
	GetHTTPStatus() int
	GetMessage() string
	GetError() error
	GetCode() string
}

type RunTimeError struct {
	ErrorText string
}

func (p RunTimeError) Error() string {
	return p.ErrorText
}
