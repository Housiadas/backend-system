package grpc

import (
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"

	userV1 "github.com/Housiadas/backend-system/gen/go/github.com/Housiadas/backend-system/gen/user/v1"
	"github.com/Housiadas/backend-system/internal/app/usecase/product_usecase"
	"github.com/Housiadas/backend-system/internal/app/usecase/system_usecase"
	"github.com/Housiadas/backend-system/internal/app/usecase/transaction_usecase"
	"github.com/Housiadas/backend-system/internal/app/usecase/user_usecase"
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
	User    *user_usecase.App
	Product *product_usecase.App
	System  *system_usecase.App
	Tx      *transaction_usecase.App
}

// Business represents the core internal layer.
type Business struct {
	User    *usercore.Core
	Product *productcore.Core
}

type Config struct {
	ServiceName string
	Build       string
	DB          *sqlx.DB
	Log         *logger.Logger
	Tracer      trace.Tracer
	UserBus     *usercore.Core
	ProductBus  *productcore.Core
}

func New(cfg Config) *Server {
	return &Server{
		ServiceName: cfg.ServiceName,
		Build:       cfg.Build,
		DB:          cfg.DB,
		Log:         cfg.Log,
		Tracer:      cfg.Tracer,
		App: App{
			User:    user_usecase.NewApp(cfg.UserBus),
			Product: product_usecase.NewApp(cfg.ProductBus),
			System:  system_usecase.NewApp(cfg.Build, cfg.Log, cfg.DB),
			Tx:      transaction_usecase.NewApp(cfg.UserBus, cfg.ProductBus),
		},
		Business: Business{
			User:    cfg.UserBus,
			Product: cfg.ProductBus,
		},
	}
}
