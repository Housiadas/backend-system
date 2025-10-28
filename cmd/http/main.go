package main

import (
	"context"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	_ "github.com/Housiadas/backend-system/docs"
	"github.com/Housiadas/backend-system/internal/app/handlers"
	"github.com/Housiadas/backend-system/internal/app/repository/audit_repo"
	"github.com/Housiadas/backend-system/internal/app/repository/product_repo"
	"github.com/Housiadas/backend-system/internal/app/repository/user_repo"
	ctxPck "github.com/Housiadas/backend-system/internal/common/context"
	"github.com/Housiadas/backend-system/internal/config"
	"github.com/Housiadas/backend-system/internal/core/service/auditcore"
	"github.com/Housiadas/backend-system/internal/core/service/authcore"
	"github.com/Housiadas/backend-system/internal/core/service/productcore"
	"github.com/Housiadas/backend-system/internal/core/service/usercore"
	"github.com/Housiadas/backend-system/pkg/debug"
	"github.com/Housiadas/backend-system/pkg/kafka"
	"github.com/Housiadas/backend-system/pkg/keystore"
	"github.com/Housiadas/backend-system/pkg/logger"
	"github.com/Housiadas/backend-system/pkg/otel"
	"github.com/Housiadas/backend-system/pkg/pgsql"
)

var build = "develop"

// @title Backend System
// @description This is a backend system.
//
// @contact.name	API Support
// @contact.url		http://www.swagger.io/support
// @contact.email	support@swagger.io
//
// @license.name	Apache 2.0
// @license.url		http://www.apache.org/licenses/LICENSE-2.0.html
//
//	@query.collection.format multi
//
// @externalDocs.description  OpenAPI
//
// @externalDocs.url	https://swagger.io/resources/open-api/
// @host				localhost:4000
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
		return ctxPck.GetRequestID(ctx)
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
	log.Info(ctx, "starting usecase", "version", cfg.Version.Build)
	defer log.Info(ctx, "shutdown complete")

	log.BuildInfo(ctx)
	expvar.NewString("build").Set(cfg.Version.Build)

	// -------------------------------------------------------------------------
	// Initialize Database
	// -------------------------------------------------------------------------
	log.Info(ctx, "startup", "status", "initializing database", "host port", cfg.DB.Host)
	db, err := pgsql.Open(pgsql.Config{
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
	// Initialize Kafka Producer
	// -------------------------------------------------------------------------
	log.Info(ctx, "startup", "status", "initializing kafka support")

	producer, err := kafka.NewProducer(kafka.ProducerConfig{
		Brokers:          cfg.Kafka.Brokers,
		LogLevel:         cfg.Kafka.LogLevel,
		AddressFamily:    cfg.Kafka.AddressFamily,
		MaxMessageBytes:  cfg.Kafka.MaxMessageBytes,
		SecurityProtocol: cfg.Kafka.SecurityProtocol,
	})
	if err != nil {
		return fmt.Errorf("creating kafka producer: %w", err)
	}
	defer producer.Close()

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
	// Core Services
	// -------------------------------------------------------------------------
	log.Info(ctx, "startup", "status", "initializing internal layer")

	auditCore := auditcore.NewCore(log, audit_repo.NewStore(log, db))
	userCore := usercore.NewCore(log, user_repo.NewStore(log, db))
	productCore := productcore.NewCore(log, userCore, product_repo.NewStore(log, db))

	// Load the private keys files from disk. We can assume some system api like
	// Vault has created these files already. How that happens is not our concern.
	ks := keystore.New()
	if err := ks.LoadRSAKeys(os.DirFS(cfg.Auth.KeysFolder)); err != nil {
		return fmt.Errorf("reading keys: %w", err)
	}
	authCore := authcore.New(authcore.Config{
		Log:       log,
		DB:        db,
		KeyLookup: ks,
		Userbus:   userCore,
	})

	// -------------------------------------------------------------------------
	// Start Debug Http Core
	// -------------------------------------------------------------------------
	go func() {
		log.Info(ctx, "startup", "status", "Debug grpc starting", "host", cfg.Http.Debug)

		if err := http.ListenAndServe(cfg.Http.Debug, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.Http.Debug, "msg", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Start API Http Server
	// -------------------------------------------------------------------------
	log.Info(ctx, "startup", "status", "API grpc starting")

	// Initialize handlers
	h := handlers.New(handlers.Config{
		ServiceName: cfg.App.Name,
		Build:       build,
		Cors:        cfg.Cors,
		DB:          db,
		Log:         log,
		Tracer:      tracer,
		AuditCore:   auditCore,
		AuthCore:    authCore,
		UserCore:    userCore,
		ProductCore: productCore,
	})

	api := http.Server{
		Addr:         cfg.Http.Api,
		Handler:      h.Routes(),
		ReadTimeout:  cfg.Http.ReadTimeout,
		WriteTimeout: cfg.Http.WriteTimeout,
		IdleTimeout:  cfg.Http.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info(ctx, "startup", "status", "API grpc started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Shutdown
	select {
	case err := <-serverErrors:
		return fmt.Errorf("grpc error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Http.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop grpc gracefully: %w", err)
		}
	}

	return nil
}
