// Package dbtest contains supporting code for running tests that hit the DB.
package dbtest

import (
	"bytes"
	"context"
	"fmt"
	ctxPck "github.com/Housiadas/backend-system/business/sys/context"
	"math/rand"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Housiadas/backend-system/business/data/sqldb"
	"github.com/Housiadas/backend-system/business/domain/productbus"
	"github.com/Housiadas/backend-system/business/domain/productbus/stores/productdb"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/domain/userbus/stores/userdb"
	"github.com/Housiadas/backend-system/pkg/docker"
	"github.com/Housiadas/backend-system/pkg/logger"
	"github.com/Housiadas/backend-system/pkg/otel"
)

const (
	PostgresImage         = "postgres:15.4"
	PostgresContainerName = "db-container"

	DBUser     = "housi"
	DBPassword = "secret123"
	DBName     = "housi_db"
	DBPort     = "5432"
)

var dbTestURL = "postgres://housi:secret123@localhost:5432/%s?sslmode=disable"

// BusDomain represents all the business domain apis needed for testing.
type BusDomain struct {
	User    *userbus.Business
	Product *productbus.Business
}

func newBusDomains(log *logger.Logger, db *sqlx.DB) BusDomain {
	userBus := userbus.NewBusiness(log, userdb.NewStore(log, db))
	productBus := productbus.NewBusiness(log, userBus, productdb.NewStore(log, db))

	return BusDomain{
		User:    userBus,
		Product: productBus,
	}
}

// =============================================================================

// Database owns state for running and shutting down tests.
type Database struct {
	DB        *sqlx.DB
	Log       *logger.Logger
	BusDomain BusDomain
}

// NewDatabase creates a new test database inside the database that was started
// to handle testing. The database is migrated to the current version and
// a connection pool is provided with business domain packages.
func NewDatabase(t *testing.T, testName string, migrationsPath string) *Database {

	dockerArgs := []string{
		"-e", "POSTGRES_DB=housi_db",
		"-e", "POSTGRES_USER=housi",
		"-e", "POSTGRES_PASSWORD=secret123",
	}
	appArgs := []string{"-c", "log_statement=all"}

	c, err := docker.StartContainer(PostgresImage, PostgresContainerName, DBPort, dockerArgs, appArgs)
	if err != nil {
		t.Fatalf("[TEST]: Starting database: %v", err)
	}

	t.Logf("Name    : %s\n", c.Name)
	t.Logf("Host: %s\n", c.HostPort)

	dbM, err := sqldb.Open(sqldb.Config{
		User:       DBUser,
		Password:   DBPassword,
		Host:       c.HostPort,
		Name:       DBName,
		DisableTLS: true,
	})
	if err != nil {
		t.Fatalf("[TEST]: Opening database connection: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := sqldb.StatusCheck(ctx, dbM); err != nil {
		t.Fatalf("[TEST]: status check database: %v", err)
	}

	// -------------------------------------------------------------------------

	const letterBytes = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 4)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	dbName := string(b)

	t.Logf("Create Database: %s\n", dbName)
	if _, err := dbM.ExecContext(context.Background(), "CREATE DATABASE "+dbName); err != nil {
		t.Fatalf("[TEST]: creating database %s: %v", dbName, err)
	}

	// -------------------------------------------------------------------------

	db, err := sqldb.Open(sqldb.Config{
		User:       DBUser,
		Password:   DBPassword,
		Host:       c.HostPort,
		Name:       dbName,
		DisableTLS: true,
	})
	if err != nil {
		t.Fatalf("[TEST]: Opening database connection: %v", err)
	}

	// -------------------------------------------------------------------------
	t.Logf("[TEST]: migrate Database UP %s\n", dbName)

	err = migration(fmt.Sprintf(dbTestURL, dbName), migrationsPath)
	if err != nil {
		t.Fatalf("[TEST]: Migrating error: %s", err)
	}

	// -------------------------------------------------------------------------

	var buf bytes.Buffer
	traceIDfn := func(context.Context) string { return otel.GetTraceID(ctx) }
	requestIDfn := func(context.Context) string { return ctxPck.GetRequestID(ctx) }
	log := logger.New(&buf, logger.LevelInfo, "TEST", traceIDfn, requestIDfn)

	// -------------------------------------------------------------------------

	// should be invoked when the caller is done with the database.
	t.Cleanup(func() {
		t.Helper()

		t.Logf("[TEST]: Drop Database: %s\n", dbName)
		if _, err := dbM.ExecContext(context.Background(), "DROP DATABASE "+dbName+" WITH (force)"); err != nil {
			t.Fatalf("[TEST]: dropping database %s: %v", dbName, err)
		}

		db.Close()
		dbM.Close()

		t.Logf("******************** LOGS (%s) ********************\n\n", testName)
		t.Log(buf.String())
		t.Logf("******************** LOGS (%s) ********************\n", testName)
	})

	return &Database{
		DB:        db,
		Log:       log,
		BusDomain: newBusDomains(log, db),
	}
}
