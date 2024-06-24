// Package dbtest contains supporting code for running tests that hit the DB.
package dbtest

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Housiadas/backend-system/business/data/sqldb"
	"github.com/Housiadas/backend-system/business/domain/productbus"
	"github.com/Housiadas/backend-system/business/domain/productbus/stores/productdb"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/domain/userbus/stores/userdb"
	"github.com/Housiadas/backend-system/business/web"
	"github.com/Housiadas/backend-system/foundation/docker"
	"github.com/Housiadas/backend-system/foundation/logger"
)

const (
	PostgresImage         = "postgres:15.4"
	PostgresContainerName = "db-container"

	DBUser     = "housi"
	DBPassword = "secret123"
	DBName     = "housi_db"
	DBPort     = "5432"
)

var migrateDbUrl = "postgres://housi:secret123@localhost:5432/%s?sslmode=disable"

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
func NewDatabase(t *testing.T, testName string) *Database {

	dockerArgs := []string{
		"-e", "POSTGRES_DB=housi_db",
		"-e", "POSTGRES_USER=housi",
		"-e", "POSTGRES_PASSWORD=secret123",
	}
	appArgs := []string{"-c", "log_statement=all"}

	c, err := docker.StartContainer(PostgresImage, PostgresContainerName, DBPort, "5432", dockerArgs, appArgs)
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

	err = migration(fmt.Sprintf(migrateDbUrl, dbName))
	if err != nil {
		t.Fatalf("[TEST]: Migrating error: %s", err)
	}

	// -------------------------------------------------------------------------

	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelInfo, "TEST", func(context.Context) string { return web.GetTraceID(ctx) })

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

// =============================================================================

// StringPointer is a helper to get a *string from a string. It is in the tests
// package because we normally don't want to deal with pointers to basic types,
// but it's useful in some tests.
func StringPointer(s string) *string {
	return &s
}

// IntPointer is a helper to get a *int from an int. It is in the tests package
// because we normally don't want to deal with pointers to basic types, but it's
// useful in some tests.
func IntPointer(i int) *int {
	return &i
}

// FloatPointer is a helper to get a *float64 from a float64. It is in the tests
// package because we normally don't want to deal with pointers to basic types,
// but it's useful in some tests.
func FloatPointer(f float64) *float64 {
	return &f
}

// BoolPointer is a helper to get a *bool from a bool. It is in the tests package
// because we normally don't want to deal with pointers to basic types, but it's
// useful in some tests.
func BoolPointer(b bool) *bool {
	return &b
}

// UserNamePointer is a helper to get a *Name from a string. It's in the tests
// package because we normally don't want to deal with pointers to basic types,
// but it's useful in some tests.
func UserNamePointer(value string) *userbus.Name {
	name := userbus.Names.MustParse(value)
	return &name
}

// ProductNamePointer is a helper to get a *Name from a string. It's in the tests
// package because we normally don't want to deal with pointers to basic types,
// but it's useful in some tests.
func ProductNamePointer(value string) *productbus.Name {
	name := productbus.Names.MustParse(value)
	return &name
}
