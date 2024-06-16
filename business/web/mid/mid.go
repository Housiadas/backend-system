// Package mid provides app level mid support.
package mid

import (
	"github.com/Housiadas/backend-system/business/auth"
	"github.com/Housiadas/backend-system/business/domain/productbus"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/foundation/logger"
)

type Mid struct {
	Auth *auth.Auth
	Bus  Business
	Log  *logger.Logger
}

type Business struct {
	Auth    *auth.Auth
	User    *userbus.Business
	Product *productbus.Business
}

func New(a *auth.Auth, b Business, l *logger.Logger) *Mid {
	return &Mid{
		Auth: a,
		Bus:  b,
		Log:  l,
	}
}
