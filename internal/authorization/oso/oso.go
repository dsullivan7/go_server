package oso

import (
	"embed"
	"fmt"
	"reflect"

	"go_server/internal/authorization"
	"go_server/internal/models"

	"github.com/osohq/go-oso"
)

//go:embed user.polar
var config embed.FS //nolint

type Authorization struct {
	oso oso.Oso
}

func NewAuthorization(o oso.Oso) authorization.Authorization {
	return &Authorization{
		oso: o,
	}
}

func (athz *Authorization) Authorize(actor interface{}, action interface{}, resource interface{}) error {
	err := athz.oso.Authorize(actor, action, resource)
	if err != nil {
		return fmt.Errorf("error in authorization: %w", err)
	}

	return nil
}

func (athz *Authorization) Init() error {
	errUser := athz.oso.RegisterClass(reflect.TypeOf(models.User{}), nil)

	if errUser != nil {
		return fmt.Errorf("error in authorization models: %w", errUser)
	}

	userAuthz, errUserAuthz := config.ReadFile("user.polar")

	if errUserAuthz != nil {
		return fmt.Errorf("error reading authorization config: %w", errUserAuthz)
	}

	errUserLoad := athz.oso.LoadString(string(userAuthz))

	if errUserLoad != nil {
		return fmt.Errorf("error loading authorization: %w", errUserLoad)
	}

	return nil
}
