package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Housiadas/backend-system/app/grpc/server"
	"github.com/Housiadas/backend-system/business/sys/dbtest"
	"github.com/Housiadas/backend-system/foundation/otel"
)

// StartTest initialized the system to run a test.
func StartTest(t *testing.T, testName string) (*Test, error) {
	db := dbtest.NewDatabase(t, testName, "file://../../../../business/data/migrations")

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

	// Initialize handler
	s := server.New(server.Config{
		ServiceName: "Test Service Name",
		Build:       "Test",
		DB:          db.DB,
		Log:         db.Log,
		Tracer:      tracer,
		UserBus:     db.BusDomain.User,
		ProductBus:  db.BusDomain.Product,
	})

	return New(db, s), nil
}
