package authorization

type Authorization interface {
	Authorize(actor interface{}, action interface{}, resource interface{}) error
}
