package grpc

import (
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"

	userV1 "github.com/Housiadas/backend-system/gen/go/github.com/Housiadas/backend-system/gen/user/v1"
	"github.com/Housiadas/backend-system/internal/app/service/productapp"
	"github.com/Housiadas/backend-system/internal/app/service/systemapp"
	"github.com/Housiadas/backend-system/internal/app/service/tranapp"
	"github.com/Housiadas/backend-system/internal/app/service/userapp"
	"github.com/Housiadas/backend-system/internal/core/service/productcore"
	"github.com/Housiadas/backend-system/internal/core/service/usercore"
	"github.com/Housiadas/backend-system/pkg/logger"
)

type Server struct {
	userV1.UnimplementedUserServiceServer
	ServiceName string
	Build       string
	DB          *sqlx.DB
	Log         *logger.Logger
	Tracer      trace.Tracer
	App         App
	Business    Business
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
	User    *usercore.Service
	Product *productcore.Business
}

type Config struct {
	ServiceName string
	Build       string
	DB          *sqlx.DB
	Log         *logger.Logger
	Tracer      trace.Tracer
	UserBus     *usercore.Service
	ProductBus  *productcore.Business
}

func New(cfg Config) *Server {
	return &Server{
		ServiceName: cfg.ServiceName,
		Build:       cfg.Build,
		DB:          cfg.DB,
		Log:         cfg.Log,
		Tracer:      cfg.Tracer,
		App: App{
			User:    userapp.NewApp(cfg.UserBus),
			Product: productapp.NewApp(cfg.ProductBus),
			System:  systemapp.NewApp(cfg.Build, cfg.Log, cfg.DB),
			Tx:      tranapp.NewApp(cfg.UserBus, cfg.ProductBus),
		},
		Business: Business{
			User:    cfg.UserBus,
			Product: cfg.ProductBus,
		},
	}
}
