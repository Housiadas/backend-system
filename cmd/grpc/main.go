package main

import (
	"context"
	"expvar"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/Housiadas/backend-system/internal/app/grpc"
	"github.com/Housiadas/backend-system/internal/app/repository/productrepo"
	"github.com/Housiadas/backend-system/internal/app/repository/userrepo"
	"github.com/Housiadas/backend-system/internal/config"
	"github.com/Housiadas/backend-system/internal/core/service/productcore"
	"github.com/Housiadas/backend-system/internal/core/service/usercore"
	"github.com/Housiadas/backend-system/pkg/logger"
	"github.com/Housiadas/backend-system/pkg/otel"
	"github.com/Housiadas/backend-system/pkg/sqldb"
)

var build = "develop"

func main() {
	// -------------------------------------------------------------------------
	// Initialize Configuration
	// -------------------------------------------------------------------------
	cfg, err := config.LoadConfig("../../")
	if err != nil {
		fmt.Errorf("parsing config: %w", err)
		os.Exit(1)
	}
	cfg.Version = config.Version{
		Build: build,
		Desc:  "gRPC",
	}

	// -------------------------------------------------------------------------
	// Initialize Logger
	// -------------------------------------------------------------------------
	var log *logger.Logger
	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
			// r.Attributes, contains all the necessary information for the alert
		},
	}
	traceIDFn := func(ctx context.Context) string {
		return otel.GetTraceID(ctx)
	}
	requestIDFn := func(ctx context.Context) string {
		return ""
	}
	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "gRPC", traceIDFn, requestIDFn, events)

	// -------------------------------------------------------------------------
	// Run the application
	// -------------------------------------------------------------------------
	ctx := context.Background()
	if err := run(ctx, cfg, log); err != nil {
		log.Error(ctx, "error during startup", "msg", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, cfg config.Config, log *logger.Logger) error {
	// -------------------------------------------------------------------------
	// GOMAXPROCS
	// -------------------------------------------------------------------------
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// App Starting
	// -------------------------------------------------------------------------
	log.Info(ctx, "starting service", "version", cfg.Version.Build)
	defer log.Info(ctx, "shutdown complete")

	log.BuildInfo(ctx)
	expvar.NewString("build").Set(cfg.Version.Build)

	// -------------------------------------------------------------------------
	// Initialize Database
	// -------------------------------------------------------------------------
	log.Info(ctx, "startup", "status", "initializing database", "host port", cfg.DB.Host)
	db, err := sqldb.Open(sqldb.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("connecting to repository: %w", err)
	}
	defer db.Close()

	// -------------------------------------------------------------------------
	// Start Tracing Support
	// -------------------------------------------------------------------------
	log.Info(ctx, "startup", "status", "initializing tracing support")

	traceProvider, teardown, err := otel.InitTracing(otel.Config{
		Log:         log,
		ServiceName: cfg.App.Name,
		Host:        cfg.Tempo.Host,
		ExcludedRoutes: map[string]struct{}{
			"/v1/liveness":  {},
			"/v1/readiness": {},
		},
		Probability: cfg.Tempo.Probability,
	})
	if err != nil {
		return fmt.Errorf("starting tracing: %w", err)
	}

	defer teardown(context.Background())

	tracer := traceProvider.Tracer(cfg.App.Name)

	// -------------------------------------------------------------------------
	// Build Core Services
	// -------------------------------------------------------------------------
	log.Info(ctx, "startup", "status", "initializing internal layer")

	userBus := usercore.NewBusiness(log, userrepo.NewStore(log, db))
	productBus := productcore.NewBusiness(log, userBus, productrepo.NewStore(log, db))

	// -------------------------------------------------------------------------
	// Start Grpc Server
	// -------------------------------------------------------------------------
	s := grpc.New(grpc.Config{
		ServiceName: cfg.App.Name,
		Build:       build,
		DB:          db,
		Log:         log,
		Tracer:      tracer,
		UserBus:     userBus,
		ProductBus:  productBus,
	})

	// Register gRPC service
	grpcServer := s.Registrar()

	listener, err := net.Listen("tcp", cfg.Grpc.Api)
	if err != nil {
		log.Error(ctx, "failed to listen", "msg", err)
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	go func() {
		log.Info(ctx, "start gRPC grpc", "address", listener.Addr().String())
		if err := grpcServer.Serve(listener); err != nil {
			log.Error(ctx, "Failed to serve gRPC grpc", err)
		}
	}()

	<-ctx.Done()
	grpcServer.GracefulStop()
	return nil
}
