package oso

import (
  "github.com/osohq/go-oso"
  "go_server/internal/authorization"
)

type Authorization struct {
	oso   oso.Oso
}

func NewAuthorization(o oso.Oso) authorization.Authorization {
	return &Authorization{
		oso:   o,
	}
}
