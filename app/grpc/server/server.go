package server

import (
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"

	"github.com/Housiadas/backend-system/app/domain/productapp"
	"github.com/Housiadas/backend-system/app/domain/systemapp"
	"github.com/Housiadas/backend-system/app/domain/tranapp"
	"github.com/Housiadas/backend-system/app/domain/userapp"
	"github.com/Housiadas/backend-system/business/domain/authbus"
	"github.com/Housiadas/backend-system/business/domain/productbus"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/foundation/logger"
)

type Server struct {
	ServiceName string
	Build       string
	DB          *sqlx.DB
	Log         *logger.Logger
	Tracer      trace.Tracer
	App         App
	Business    Business
}

// App represents the core app layer
type App struct {
	User    *userapp.App
	Product *productapp.App
	System  *systemapp.App
	Tx      *tranapp.App
}

// Business represents the core business layer.
type Business struct {
	Auth    *authbus.Auth
	User    *userbus.Business
	Product *productbus.Business
}
