package handler

import (
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"

	"github.com/Housiadas/backend-system/app/domain/productapp"
	"github.com/Housiadas/backend-system/app/domain/systemapp"
	"github.com/Housiadas/backend-system/app/domain/tranapp"
	"github.com/Housiadas/backend-system/app/domain/userapp"
	"github.com/Housiadas/backend-system/internal/config"
	"github.com/Housiadas/backend-system/internal/data/sqldb"
	"github.com/Housiadas/backend-system/internal/domain/authbus"
	"github.com/Housiadas/backend-system/internal/domain/productbus"
	"github.com/Housiadas/backend-system/internal/domain/userbus"
	"github.com/Housiadas/backend-system/internal/sys/middleware"
	"github.com/Housiadas/backend-system/pkg/logger"
	"github.com/Housiadas/backend-system/pkg/web"
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
	Middleware *middleware.Middleware
	Res        *web.Respond
}

// App represents the core app layer
type App struct {
	User    *userapp.App
	Product *productapp.App
	System  *systemapp.App
	Tx      *tranapp.App
}

// Business represents the core internal layer.
type Business struct {
	Auth    *authbus.Auth
	User    *userbus.Business
	Product *productbus.Business
}

// Config represents the configuration for the handler.
type Config struct {
	ServiceName string
	Build       string
	Cors        config.CorsSettings
	DB          *sqlx.DB
	Log         *logger.Logger
	Tracer      trace.Tracer
	AuthBus     *authbus.Auth
	UserBus     *userbus.Business
	ProductBus  *productbus.Business
}

func New(cfg Config) *Handler {
	return &Handler{
		ServiceName: cfg.ServiceName,
		Build:       cfg.Build,
		Cors:        cfg.Cors,
		DB:          cfg.DB,
		Log:         cfg.Log,
		Tracer:      cfg.Tracer,
		Web: Web{
			Middleware: middleware.New(middleware.Config{
				Log:     cfg.Log,
				Tracer:  cfg.Tracer,
				Tx:      sqldb.NewBeginner(cfg.DB),
				Auth:    cfg.AuthBus,
				User:    cfg.UserBus,
				Product: cfg.ProductBus,
			}),
			Res: web.NewRespond(cfg.Log),
		},
		App: App{
			User:    userapp.NewAppWithAuth(cfg.UserBus, cfg.AuthBus),
			Product: productapp.NewApp(cfg.ProductBus),
			System:  systemapp.NewApp(cfg.Build, cfg.Log, cfg.DB),
			Tx:      tranapp.NewApp(cfg.UserBus, cfg.ProductBus),
		},
		Business: Business{
			Auth:    cfg.AuthBus,
			User:    cfg.UserBus,
			Product: cfg.ProductBus,
		},
	}
}
