package middleware

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Housiadas/backend-system/internal/data/sqldb"
	"github.com/Housiadas/backend-system/pkg/errs"
)

// BeginCommitRollback starts a transaction for the domain call.
func (m *Middleware) BeginCommitRollback() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			hasCommitted := false

			m.Log.Info(ctx, "BEGIN TRANSACTION")
			tx, err := m.Tx.Begin()
			if err != nil {
				err := errs.Newf(errs.Internal, "BEGIN TRANSACTION: %s", err)
				m.Log.Error(ctx, "transaction middleware", err)
				m.Error(w, err, http.StatusInternalServerError)
				return
			}

			defer func() {
				if !hasCommitted {
					m.Log.Info(ctx, "ROLLBACK TRANSACTION")
				}

				if err := tx.Rollback(); err != nil {
					if errors.Is(err, sql.ErrTxDone) {
						return
					}
					m.Log.Error(ctx, "ROLLBACK TRANSACTION", "ERROR", err)
				}
			}()

			ctx = sqldb.SetTran(ctx, tx)

			// Create a response recorder to capture the response
			rec := &ResponseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rec, r.WithContext(ctx))

			// Access the recorded response
			// Check if we can commit transaction
			if rec.statusCode >= 400 {
				m.Log.Info(ctx, "TRANSACTION FAILED, WILL ROLLBACK")
				return
			}

			m.Log.Info(ctx, "COMMIT TRANSACTION")
			if err := tx.Commit(); err != nil {
				err := errs.Newf(errs.Internal, "COMMIT TRANSACTION: %s", err)
				m.Log.Error(ctx, "transaction middleware", err)
				m.Error(w, err, http.StatusInternalServerError)
				return
			}

			hasCommitted = true
		})
	}
}
