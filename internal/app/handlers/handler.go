package handlers

import (
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"

	"github.com/Housiadas/backend-system/internal/app/middleware"
	"github.com/Housiadas/backend-system/internal/app/usecase/audit_usecase"
	"github.com/Housiadas/backend-system/internal/app/usecase/product_usecase"
	"github.com/Housiadas/backend-system/internal/app/usecase/system_usecase"
	"github.com/Housiadas/backend-system/internal/app/usecase/transaction_usecase"
	"github.com/Housiadas/backend-system/internal/app/usecase/user_usecase"
	"github.com/Housiadas/backend-system/internal/config"
	"github.com/Housiadas/backend-system/internal/core/service/auditcore"
	"github.com/Housiadas/backend-system/internal/core/service/authcore"
	"github.com/Housiadas/backend-system/internal/core/service/productcore"
	"github.com/Housiadas/backend-system/internal/core/service/usercore"
	"github.com/Housiadas/backend-system/pkg/logger"
	"github.com/Housiadas/backend-system/pkg/pgsql"
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

// Web represents the set of usecase for the http.
type Web struct {
	Middleware *middleware.Middleware
	Res        *web.Respond
}

// App represents the core cli layer
type App struct {
	Audit   *audit_usecase.App
	User    *user_usecase.App
	Product *product_usecase.App
	System  *system_usecase.App
	Tx      *transaction_usecase.App
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
				Tx:      pgsql.NewBeginner(cfg.DB),
				Auth:    cfg.AuthCore,
				User:    cfg.UserCore,
				Product: cfg.ProductCore,
			}),
			Res: web.NewRespond(cfg.Log),
		},
		App: App{
			Audit:   audit_usecase.NewApp(cfg.AuditCore),
			User:    user_usecase.NewAppWithAuth(cfg.UserCore, cfg.AuthCore),
			Product: product_usecase.NewApp(cfg.ProductCore),
			System:  system_usecase.NewApp(cfg.Build, cfg.Log, cfg.DB),
			Tx:      transaction_usecase.NewApp(cfg.UserCore, cfg.ProductCore),
		},
		Core: Core{
			Audit:   cfg.AuditCore,
			Auth:    cfg.AuthCore,
			User:    cfg.UserCore,
			Product: cfg.ProductCore,
		},
	}
}
