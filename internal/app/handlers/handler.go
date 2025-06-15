package handlers

import (
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"

	"github.com/Housiadas/backend-system/internal/app/middleware"
	"github.com/Housiadas/backend-system/internal/app/service/auditapp"
	"github.com/Housiadas/backend-system/internal/app/service/productapp"
	"github.com/Housiadas/backend-system/internal/app/service/systemapp"
	"github.com/Housiadas/backend-system/internal/app/service/tranapp"
	"github.com/Housiadas/backend-system/internal/app/service/userapp"
	"github.com/Housiadas/backend-system/internal/config"
	"github.com/Housiadas/backend-system/internal/core/service/auditcore"
	"github.com/Housiadas/backend-system/internal/core/service/authcore"
	"github.com/Housiadas/backend-system/internal/core/service/productcore"
	"github.com/Housiadas/backend-system/internal/core/service/usercore"
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
	Core        Core
}

// Web represents the set of service for the http.
type Web struct {
	Middleware *middleware.Middleware
	Res        *web.Respond
}

// App represents the core cli layer
type App struct {
	Audit   *auditapp.App
	User    *userapp.App
	Product *productapp.App
	System  *systemapp.App
	Tx      *tranapp.App
}

// Core represents the core internal layer.
type Core struct {
	Auth    *authcore.Auth
	Audit   *auditcore.Core
	User    *usercore.Core
	Product *productcore.Core
}

// Config represents the configuration for the handlers.
type Config struct {
	ServiceName string
	Build       string
	Cors        config.CorsSettings
	DB          *sqlx.DB
	Log         *logger.Logger
	Tracer      trace.Tracer
	AuditCore   *auditcore.Core
	AuthCore    *authcore.Auth
	UserCore    *usercore.Core
	ProductCore *productcore.Core
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
				Auth:    cfg.AuthCore,
				User:    cfg.UserCore,
				Product: cfg.ProductCore,
			}),
			Res: web.NewRespond(cfg.Log),
		},
		App: App{
			Audit:   auditapp.NewApp(cfg.AuditCore),
			User:    userapp.NewAppWithAuth(cfg.UserCore, cfg.AuthCore),
			Product: productapp.NewApp(cfg.ProductCore),
			System:  systemapp.NewApp(cfg.Build, cfg.Log, cfg.DB),
			Tx:      tranapp.NewApp(cfg.UserCore, cfg.ProductCore),
		},
		Core: Core{
			Audit:   cfg.AuditCore,
			Auth:    cfg.AuthCore,
			User:    cfg.UserCore,
			Product: cfg.ProductCore,
		},
	}
}
