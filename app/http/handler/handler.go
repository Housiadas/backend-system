package handler

import (
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"

	"github.com/Housiadas/backend-system/app/domain/productapp"
	"github.com/Housiadas/backend-system/app/domain/systemapp"
	"github.com/Housiadas/backend-system/app/domain/tranapp"
	"github.com/Housiadas/backend-system/app/domain/userapp"
	"github.com/Housiadas/backend-system/business/config"
	"github.com/Housiadas/backend-system/business/domain/authbus"
	"github.com/Housiadas/backend-system/business/domain/productbus"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/web"
	"github.com/Housiadas/backend-system/business/web/mid"
	"github.com/Housiadas/backend-system/foundation/logger"
)

// Handler contains all the mandatory systems required by handlers.
type Handler struct {
	ServiceName string
	Build       string
	Cors        config.CorsSettings
	DB          *sqlx.DB
	Log         *logger.Logger
	Tracer      trace.Tracer
	Web         Web
	App         App
	Business    Business
}

// Web represents the set of services for the http.
type Web struct {
	Mid *mid.Mid
	Res *web.Respond
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
