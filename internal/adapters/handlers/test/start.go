package test

import (
	"context"
	"fmt"
	"github.com/Housiadas/backend-system/internal/adapters/repository/userrepository"
	"testing"

	"github.com/Housiadas/backend-system/internal/adapters/handlers"
	"github.com/Housiadas/backend-system/internal/common/dbtest"
	cfg "github.com/Housiadas/backend-system/internal/config"
	"github.com/Housiadas/backend-system/internal/core/service/authbus"
	"github.com/Housiadas/backend-system/internal/core/service/userservice"
	"github.com/Housiadas/backend-system/pkg/otel"
)

// StartTest initialized the system to run a test.
func StartTest(t *testing.T, testName string) (*Test, error) {
	db := dbtest.NewDatabase(t, testName)

	// auth
	auth := authbus.New(authbus.Config{
		Log:       db.Log,
		DB:        db.DB,
		KeyLookup: &KeyStore{},
		Userbus:   userservice.NewBusiness(db.Log, userrepository.NewStore(db.Log, db.DB)),
	})

	// tracer
	traceProvider, teardown, err := otel.InitTracing(otel.Config{
		Log:         db.Log,
		ServiceName: "Service Name",
		Host:        "Test host",
		ExcludedRoutes: map[string]struct{}{
			"/v1/liveness":  {},
			"/v1/readiness": {},
		},
		Probability: 0.5,
	})
	if err != nil {
		return nil, fmt.Errorf("starting tracing: %w", err)
	}

	defer teardown(context.Background())

	tracer := traceProvider.Tracer("Service Name")

	// Initialize handlers
	h := handlers.New(handlers.Config{
		ServiceName: "Test Service Name",
		Build:       "Test",
		Cors:        cfg.CorsSettings{},
		DB:          db.DB,
		Log:         db.Log,
		Tracer:      tracer,
		AuthBus:     auth,
		UserBus:     db.BusDomain.User,
		ProductBus:  db.BusDomain.Product,
	})

	return New(db, auth, h.Routes()), nil
}
