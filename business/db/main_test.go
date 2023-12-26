package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource = "postgres://housi:secret123@localhost/housi_db?sslmode=disable"
)

var testStore *Queries

func TestMain(m *testing.M) {
	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testStore = New(connPool)
	os.Exit(m.Run())
}
