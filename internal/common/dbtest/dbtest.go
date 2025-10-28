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

	ctxPck "github.com/Housiadas/backend-system/internal/common/context"
	"github.com/Housiadas/backend-system/pkg/docker"
	"github.com/Housiadas/backend-system/pkg/logger"
	"github.com/Housiadas/backend-system/pkg/otel"
	"github.com/Housiadas/backend-system/pkg/pgsql"
)

const (
	PostgresImage         = "postgres:15.4"
	PostgresContainerName = "repository-container"

	DBUser     = "housi"
	DBPassword = "secret123"
	DBName     = "housi_db"
	DBPort     = "5432"
)

var dbTestURL = "postgres://housi:secret123@localhost:5432/%s?sslmode=disable"

// Database owns the state for running and shutting down tests.
type Database struct {
	DB   *sqlx.DB
	Log  *logger.Logger
	Core Core
}

// New creates a new test database inside the database that was started
// to handle testing. The database is migrated to the current version, and
// a connection pool is provided with internal core packages.
func New(t *testing.T, testName string) *Database {

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

	dbM, err := pgsql.Open(pgsql.Config{
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

	if err := pgsql.StatusCheck(ctx, dbM); err != nil {
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

	db, err := pgsql.Open(pgsql.Config{
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

	err = migration(fmt.Sprintf(dbTestURL, dbName))
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
		DB:   db,
		Log:  log,
		Core: newCore(log, db),
	}
}
