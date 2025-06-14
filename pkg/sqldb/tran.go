package sqldb

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ctxKey string

const (
	tranKey ctxKey = "tranKey"
)

// Beginner represents a value that can begin a transaction.
type Beginner interface {
	Begin() (CommitRollbacker, error)
}

// CommitRollbacker represents a value that can commit or rollback a transaction.
type CommitRollbacker interface {
	Commit() error
	Rollback() error
}

// =============================================================================

// DBBeginner implements the Beginner interface,
type DBBeginner struct {
	sqlxDB *sqlx.DB
}

// NewBeginner constructs a value that implements the beginner interface.
func NewBeginner(sqlxDB *sqlx.DB) *DBBeginner {
	return &DBBeginner{
		sqlxDB: sqlxDB,
	}
}

// Begin implements the Beginner interface and returns a concrete value that
// implements the CommitRollbacker interface.
func (db *DBBeginner) Begin() (CommitRollbacker, error) {
	return db.sqlxDB.Beginx()
}

// GetExtContext is a helper function that extracts the sqlx value
// from the core transactor interface for transactional use.
func GetExtContext(tx CommitRollbacker) (sqlx.ExtContext, error) {
	ec, ok := tx.(sqlx.ExtContext)
	if !ok {
		return nil, fmt.Errorf("Transactor(%T) not of a type *sql.Tx", tx)
	}

	return ec, nil
}

func SetTran(ctx context.Context, tx CommitRollbacker) context.Context {
	return context.WithValue(ctx, tranKey, tx)
}

// GetTran retrieves the value that can manage a transaction.
func GetTran(ctx context.Context) (CommitRollbacker, error) {
	v, ok := ctx.Value(tranKey).(CommitRollbacker)
	if !ok {
		return nil, errors.New("transaction not found in context")
	}

	return v, nil
}
