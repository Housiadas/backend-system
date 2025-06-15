package apitest

import (
	"context"
	"fmt"
	"testing"

	"github.com/Housiadas/backend-system/internal/app/handlers"
	"github.com/Housiadas/backend-system/internal/app/repository/userrepo"
	"github.com/Housiadas/backend-system/internal/common/dbtest"
	cfg "github.com/Housiadas/backend-system/internal/config"
	"github.com/Housiadas/backend-system/internal/core/service/authcore"
	"github.com/Housiadas/backend-system/internal/core/service/usercore"
	"github.com/Housiadas/backend-system/pkg/otel"
)

// StartTest initialized the system to run a test.
func StartTest(t *testing.T, testName string) (*Test, error) {
	db := dbtest.New(t, testName)

	// auth
	auth := authcore.New(authcore.Config{
		Log:       db.Log,
		DB:        db.DB,
		KeyLookup: &KeyStore{},
		Userbus:   usercore.NewCore(db.Log, userrepo.NewStore(db.Log, db.DB)),
	})

	// tracer
	traceProvider, teardown, err := otel.InitTracing(otel.Config{
		Log:         db.Log,
		ServiceName: "Core Name",
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

	tracer := traceProvider.Tracer("Core Name")

	// Initialize handlers
	h := handlers.New(handlers.Config{
		ServiceName: "Test Service Name",
		Build:       "Test",
		Cors:        cfg.CorsSettings{},
		DB:          db.DB,
		Log:         db.Log,
		Tracer:      tracer,
		AuthCore:    auth,
		UserCore:    db.Core.User,
		ProductCore: db.Core.Product,
	})

	return New(db, auth, h.Routes()), nil
}
