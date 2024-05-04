package main

import (
	"context"
	"expvar"
	"fmt"
	"github.com/Housiadas/backend-system/foundation/kafka"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/Housiadas/backend-system/app/api/mux"
	"github.com/Housiadas/backend-system/app/api/route"
	"github.com/Housiadas/backend-system/business/auth"
	"github.com/Housiadas/backend-system/business/config"
	"github.com/Housiadas/backend-system/business/data/sqldb"
	"github.com/Housiadas/backend-system/business/domain/productbus"
	"github.com/Housiadas/backend-system/business/domain/productbus/stores/productdb"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/domain/userbus/stores/userdb"
	del "github.com/Housiadas/backend-system/business/sys/delegate"
	"github.com/Housiadas/backend-system/foundation/debug"
	"github.com/Housiadas/backend-system/foundation/keystore"
	"github.com/Housiadas/backend-system/foundation/logger"
	"github.com/Housiadas/backend-system/foundation/web"
)

/*
	Need to figure out timeouts for http service.
*/

var build = "develop"

func main() {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}
	traceIDFn := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "API", traceIDFn, events)

	// -------------------------------------------------------------------------
	// Run the application
	// -------------------------------------------------------------------------
	ctx := context.Background()
	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "msg", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS
	// -------------------------------------------------------------------------
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration
	// -------------------------------------------------------------------------
	cfg, err := config.LoadConfig("../../")
	if err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}
	cfg.Version = config.Version{
		Build: build,
		Desc:  "API",
	}

	// -------------------------------------------------------------------------
	// App Starting
	// -------------------------------------------------------------------------
	log.Info(ctx, "starting service", "version", cfg.Version.Build)
	defer log.Info(ctx, "shutdown complete")
	expvar.NewString("build").Set(cfg.Version.Build)

	// -------------------------------------------------------------------------
	// Database Support
	// -------------------------------------------------------------------------
	log.Info(ctx, "startup", "status", "initializing database support", "hostport", cfg.DB.Host)

	db, err := sqldb.Open(sqldb.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		HostPort:     cfg.DB.Host,
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
	// Initialize authentication support
	// -------------------------------------------------------------------------
	log.Info(ctx, "startup", "status", "initializing authentication support")

	// Load the private keys files from disk. We can assume some system like
	// Vault has created these files already. How that happens is not our concern.
	ks := keystore.New()
	if err := ks.LoadRSAKeys(os.DirFS(cfg.Auth.KeysFolder)); err != nil {
		return fmt.Errorf("reading keys: %w", err)
	}

	authCfg := auth.Config{
		Log:       log,
		DB:        db,
		KeyLookup: ks,
	}

	authSrv, err := auth.New(authCfg)
	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}

	// -------------------------------------------------------------------------
	// Initialize Kafka Producer
	// -------------------------------------------------------------------------
	producer, err := kafka.NewProducer(kafka.ProducerConfig{
		Broker:           cfg.Kafka.Broker,
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

	//log.Info(ctx, "startup", "status", "initializing tracing support")
	//
	//traceProvider, err := startTracing(
	//	cfg.Tempo.ServiceName,
	//	cfg.Tempo.ReporterURI,
	//	cfg.Tempo.Probability,
	//)
	//if err != nil {
	//	return fmt.Errorf("starting tracing: %w", err)
	//}
	//
	//defer traceProvider.Shutdown(context.Background())
	//
	//tracer := traceProvider.Tracer("service")

	// -------------------------------------------------------------------------
	// Build Core APIs
	// -------------------------------------------------------------------------
	log.Info(ctx, "startup", "status", "initializing business core")

	delegate := del.New(log)
	userBus := userbus.NewCore(log, userdb.NewStore(log, db), producer)
	productBus := productbus.NewCore(log, productdb.NewStore(log, db), userBus, producer)

	// -------------------------------------------------------------------------
	// Start Debug Service
	// -------------------------------------------------------------------------
	go func() {
		log.Info(ctx, "startup", "status", "debug v1 router started", "host", cfg.Server.Debug)

		if err := http.ListenAndServe(cfg.Server.Debug, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.Server.Debug, "msg", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Start API Service
	// -------------------------------------------------------------------------
	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfgMux := mux.Config{
		Build:    build,
		Shutdown: shutdown,
		Log:      log,
		DB:       db,
		Auth:     authSrv,
		BusDomain: mux.BusDomain{
			Delegate: delegate,
			User:     userBus,
			Product:  productBus,
		},
	}

	api := http.Server{
		Addr:         cfg.Server.Api,
		Handler:      mux.WebAPI(cfgMux, route.Routes(), mux.WithCORS(cfg.Server.CorsAllowedOrigins)),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown
	// -------------------------------------------------------------------------
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Server.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
