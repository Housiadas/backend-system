package main

import (
	"context"
	"expvar"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/Housiadas/backend-system/app/grpc/server"
	"github.com/Housiadas/backend-system/business/config"
	"github.com/Housiadas/backend-system/business/data/sqldb"
	"github.com/Housiadas/backend-system/business/domain/productbus"
	"github.com/Housiadas/backend-system/business/domain/productbus/stores/productdb"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/domain/userbus/stores/userdb"
	"github.com/Housiadas/backend-system/business/web"
	"github.com/Housiadas/backend-system/foundation/logger"
	"github.com/Housiadas/backend-system/foundation/otel"
	userV1 "github.com/Housiadas/backend-system/gen/go/github.com/Housiadas/backend-system/gen/user/v1"
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
		return web.GetRequestID(ctx)
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
		return fmt.Errorf("connecting to db: %w", err)
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
	// Build Business Layer
	// -------------------------------------------------------------------------
	log.Info(ctx, "startup", "status", "initializing business layer")

	userBus := userbus.NewBusiness(log, userdb.NewStore(log, db))
	productBus := productbus.NewBusiness(log, userBus, productdb.NewStore(log, db))

	// Initialize Server Struct
	s := server.New(server.Config{
		ServiceName: cfg.App.Name,
		Build:       build,
		DB:          db,
		Log:         log,
		Tracer:      tracer,
		UserBus:     userBus,
		ProductBus:  productBus,
	})

	// -------------------------------------------------------------------------
	// Health Server
	// -------------------------------------------------------------------------
	healthServer := health.NewServer()
	go func() {
		for {
			status := healthpb.HealthCheckResponse_SERVING
			// Check if user Service is valid
			if time.Now().Second()%2 == 0 {
				status = healthpb.HealthCheckResponse_NOT_SERVING
			}

			healthServer.SetServingStatus(userV1.UserService_ServiceDesc.ServiceName, status)
			healthServer.SetServingStatus("", status)

			time.Sleep(1 * time.Second)
		}
	}()

	// -------------------------------------------------------------------------
	// Register gRPC services
	// -------------------------------------------------------------------------
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(s.GrpcInterceptor),
	)
	userV1.RegisterUserServiceServer(grpcServer, s)
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	reflection.Register(grpcServer)

	// -------------------------------------------------------------------------
	// Start Grpc Server
	// -------------------------------------------------------------------------
	listener, err := net.Listen("tcp", cfg.Grpc.Api)
	if err != nil {
		log.Error(ctx, "failed to listen", "msg", err)
	}

	log.Info(ctx, "start gRPC server", "address", listener.Addr().String())

	// todo add graceful shutdown
	return grpcServer.Serve(listener)
}
