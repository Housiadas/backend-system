package testint

import (
	"context"
	"fmt"
	"testing"

	"github.com/Housiadas/backend-system/app/api/handler"
	"github.com/Housiadas/backend-system/app/domain/productapp"
	"github.com/Housiadas/backend-system/app/domain/systemapp"
	"github.com/Housiadas/backend-system/app/domain/tranapp"
	"github.com/Housiadas/backend-system/app/domain/userapp"
	cfg "github.com/Housiadas/backend-system/business/config"
	"github.com/Housiadas/backend-system/business/data/sqldb"
	"github.com/Housiadas/backend-system/business/domain/authbus"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/domain/userbus/stores/userdb"
	"github.com/Housiadas/backend-system/business/sys/dbtest"
	"github.com/Housiadas/backend-system/business/web"
	"github.com/Housiadas/backend-system/business/web/mid"
	"github.com/Housiadas/backend-system/foundation/otel"
)

// StartTest initialized the system to run a test.
func StartTest(t *testing.T, testName string) (*Test, error) {
	db := dbtest.NewDatabase(t, testName)

	auth := authbus.New(authbus.Config{
		Log:       db.Log,
		DB:        db.DB,
		KeyLookup: &KeyStore{},
		Userbus:   userbus.NewBusiness(db.Log, userdb.NewStore(db.Log, db.DB)),
	})

	traceProvider, err := otel.InitTracing(otel.Config{
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
	defer traceProvider.Shutdown(context.Background())
	tracer := traceProvider.Tracer("Service Name")

	// Initialize handler
	h := handler.Handler{
		ServiceName: "Test Service Name",
		Build:       "Test",
		Cors:        cfg.CorsSettings{},
		DB:          db.DB,
		Log:         db.Log,
		Tracer:      tracer,
		Web: handler.Web{
			Mid: mid.New(
				mid.Business{
					Auth:    auth,
					User:    db.BusDomain.User,
					Product: db.BusDomain.Product,
				},
				db.Log,
				tracer,
				sqldb.NewBeginner(db.DB),
			),
			Res: web.NewRespond(db.Log),
		},
		App: handler.App{
			User:    userapp.NewApp(db.BusDomain.User, auth),
			Product: productapp.NewApp(db.BusDomain.Product),
			System:  systemapp.NewApp("Test version", db.Log, db.DB),
			Tx:      tranapp.NewApp(db.BusDomain.User, db.BusDomain.Product),
		},
		Business: handler.Business{
			Auth:    auth,
			User:    db.BusDomain.User,
			Product: db.BusDomain.Product,
		},
	}

	return New(db, auth, h.Routes()), nil
}
