package handlers

import (
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"

	"github.com/Housiadas/backend-system/internal/adapters/domain/productapp"
	"github.com/Housiadas/backend-system/internal/adapters/domain/systemapp"
	"github.com/Housiadas/backend-system/internal/adapters/domain/tranapp"
	"github.com/Housiadas/backend-system/internal/adapters/domain/userapp"
	"github.com/Housiadas/backend-system/internal/adapters/middleware"
	"github.com/Housiadas/backend-system/internal/config"
	"github.com/Housiadas/backend-system/internal/core/service/authbus"
	"github.com/Housiadas/backend-system/internal/core/service/productservice"
	"github.com/Housiadas/backend-system/internal/core/service/userservice"
	"github.com/Housiadas/backend-system/pkg/logger"
	"github.com/Housiadas/backend-system/pkg/sqldb"
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

// Web represents the set of service for the http.
type Web struct {
	Middleware *middleware.Middleware
	Res        *web.Respond
}

// App represents the core cli layer
type App struct {
	User    *userapp.App
	Product *productapp.App
	System  *systemapp.App
	Tx      *tranapp.App
}

// Business represents the core internal layer.
type Business struct {
	Auth    *authbus.Auth
	User    *userservice.Service
	Product *productservice.Business
}

// Config represents the configuration for the handlers.
type Config struct {
	ServiceName string
	Build       string
	Cors        config.CorsSettings
	DB          *sqlx.DB
	Log         *logger.Logger
	Tracer      trace.Tracer
	AuthBus     *authbus.Auth
	UserBus     *userservice.Service
	ProductBus  *productservice.Business
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
