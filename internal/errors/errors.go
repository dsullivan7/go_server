package errors

type HTTPError interface {
	GetHTTPStatus() int
	GetMessage() string
	GetError() error
}

type RunTimeError struct {
	ErrorText string
}

func (p RunTimeError) Error() string {
	return p.ErrorText
}
