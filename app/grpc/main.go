package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	userv1 "github.com/Housiadas/backend-system/gen/go/github.com/Housiadas/backend-system/gen/user/v1"
	"google.golang.org/grpc"
	"net"
	"os"
	"runtime"

	"github.com/Housiadas/backend-system/app/domain/productapp"
	"github.com/Housiadas/backend-system/app/domain/systemapp"
	"github.com/Housiadas/backend-system/app/domain/tranapp"
	"github.com/Housiadas/backend-system/app/domain/userapp"
	"github.com/Housiadas/backend-system/app/grpc/server"
	"github.com/Housiadas/backend-system/business/config"
	"github.com/Housiadas/backend-system/business/data/sqldb"
	"github.com/Housiadas/backend-system/business/domain/authbus"
	"github.com/Housiadas/backend-system/business/domain/productbus"
	"github.com/Housiadas/backend-system/business/domain/productbus/stores/productdb"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/domain/userbus/stores/userdb"
	"github.com/Housiadas/backend-system/business/web"
	_ "github.com/Housiadas/backend-system/docs"
	"github.com/Housiadas/backend-system/foundation/keystore"
	"github.com/Housiadas/backend-system/foundation/logger"
	"github.com/Housiadas/backend-system/foundation/otel"
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
		Desc:  "API",
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
	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "API", traceIDFn, requestIDFn, events)

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

	// Load the private keys files from disk. We can assume some system api like
	// Vault has created these files already. How that happens is not our concern.
	ks := keystore.New()
	if err := ks.LoadRSAKeys(os.DirFS(cfg.Auth.KeysFolder)); err != nil {
		return fmt.Errorf("reading keys: %w", err)
	}
	authBus := authbus.New(authbus.Config{
		Log:       log,
		DB:        db,
		KeyLookup: ks,
		Userbus:   userBus,
	})

	// -------------------------------------------------------------------------
	// Start API Server
	// -------------------------------------------------------------------------
	log.Info(ctx, "startup", "status", "API server starting")

	// Initialize handler
	s := server.Server{
		ServiceName: cfg.App.Name,
		Build:       build,
		DB:          db,
		Log:         log,
		Tracer:      tracer,
		App: server.App{
			User:    userapp.NewApp(userBus, authBus),
			Product: productapp.NewApp(productBus),
			System:  systemapp.NewApp(cfg.Version.Build, log, db),
			Tx:      tranapp.NewApp(userBus, productBus),
		},
		Business: server.Business{
			Auth:    authBus,
			User:    userBus,
			Product: productBus,
		},
	}

	gprcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(gprcLogger)
	userv1.File_user_v1_user_service_proto(grpcServer, s)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start gRPC server at %s", listener.Addr().String())

		err = grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Err(err).Msg("gRPC server failed to serve")
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown gRPC server")

		grpcServer.GracefulStop()
		log.Info().Msg("gRPC server is stopped")

		return nil
	})

	return nil
}
