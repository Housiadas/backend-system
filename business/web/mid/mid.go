// Package mid provides app level mid support.
package mid

import (
	"go.opentelemetry.io/otel/trace"

	"github.com/Housiadas/backend-system/business/domain/authbus"
	"github.com/Housiadas/backend-system/business/domain/productbus"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/foundation/logger"
)

type Mid struct {
	Bus    Business
	Log    *logger.Logger
	Tracer trace.Tracer
}

type Business struct {
	Auth    *authbus.Auth
	User    *userbus.Business
	Product *productbus.Business
}

func New(b Business, l *logger.Logger, t trace.Tracer) *Mid {
	return &Mid{
		Bus:    b,
		Log:    l,
		Tracer: t,
	}
}
