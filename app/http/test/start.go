package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Housiadas/backend-system/app/http/handler"
	cfg "github.com/Housiadas/backend-system/business/config"
	"github.com/Housiadas/backend-system/business/domain/authbus"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/domain/userbus/stores/userdb"
	"github.com/Housiadas/backend-system/business/sys/dbtest"
	"github.com/Housiadas/backend-system/foundation/otel"
)

// StartTest initialized the system to run a test.
func StartTest(t *testing.T, testName string) (*Test, error) {
	db := dbtest.NewDatabase(t, testName, "file://../../../../business/data/migrations")

	// auth
	auth := authbus.New(authbus.Config{
		Log:       db.Log,
		DB:        db.DB,
		KeyLookup: &KeyStore{},
		Userbus:   userbus.NewBusiness(db.Log, userdb.NewStore(db.Log, db.DB)),
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

	// Initialize handler
	h := handler.New(handler.Config{
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
