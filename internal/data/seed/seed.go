package seed

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/Housiadas/backend-system/internal/data/sqldb"
	"github.com/jmoiron/sqlx"
)

var (
	//go:embed seed.sql
	seedDoc string
)

// Seeder runs the seed document defined in this package against db. The queries
// are run in a transaction and rolled back if any fail.
func Seeder(ctx context.Context, db *sqlx.DB) (err error) {
	if err := sqldb.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if errTx := tx.Rollback(); errTx != nil {
			if errors.Is(errTx, sql.ErrTxDone) {
				return
			}

			err = fmt.Errorf("rollback: %w", errTx)
			return
		}
	}()

	if _, err := tx.Exec(seedDoc); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}
