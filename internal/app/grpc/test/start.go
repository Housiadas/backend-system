package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Housiadas/backend-system/internal/app/grpc"
	"github.com/Housiadas/backend-system/internal/common/dbtest"
	"github.com/Housiadas/backend-system/pkg/otel"
)

// StartTest initialized the system to run a test.
func StartTest(t *testing.T, testName string) (*Test, error) {
	db := dbtest.New(t, testName)

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
	s := grpc.New(grpc.Config{
		ServiceName: "Test Service Name",
		Build:       "Test",
		DB:          db.DB,
		Log:         db.Log,
		Tracer:      tracer,
		UserBus:     db.Core.User,
		ProductBus:  db.Core.Product,
	})

	return New(db, s), nil
}
