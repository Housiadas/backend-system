package mid

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Housiadas/backend-system/business/data/sqldb"
	"github.com/Housiadas/backend-system/business/sys/errs"
)

// BeginCommitRollback starts a transaction for the domain call.
func (m *Mid) BeginCommitRollback() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			hasCommitted := false

			m.Log.Info(ctx, "BEGIN TRANSACTION")
			tx, err := m.Tx.Begin()
			if err != nil {
				err := errs.Newf(errs.Internal, "BEGIN TRANSACTION: %s", err)
				m.Log.Error(ctx, "transaction middleware", err)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				if err := json.NewEncoder(w).Encode(err); err != nil {
					return
				}
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
			if rec.statusCode != 200 {
				m.Log.Info(ctx, "TRANSACTION STATUS FALSE")
				return
			}

			m.Log.Info(ctx, "COMMIT TRANSACTION")
			if err := tx.Commit(); err != nil {
				err := errs.Newf(errs.Internal, "COMMIT TRANSACTION: %s", err)
				m.Log.Error(ctx, "transaction middleware", err)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				if err := json.NewEncoder(w).Encode(err); err != nil {
					return
				}
			}

			hasCommitted = true
		})
	}
}
