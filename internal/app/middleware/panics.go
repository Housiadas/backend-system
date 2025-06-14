package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/Housiadas/backend-system/pkg/errs"
	"github.com/Housiadas/backend-system/pkg/metrics"
)

// Recoverer recovers from panics and converts the panic to an error,
// so it is reported in Metrics and handled in Errors.
func (m *Middleware) Recoverer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Defer a function to recover from panic
			defer func() {
				if rec := recover(); rec != nil {
					metrics.AddPanics(ctx)
					trace := debug.Stack()
					err := errs.Newf(errs.InternalOnlyLog, "PANIC [%v] TRACE[%s]", rec, string(trace))
					m.Log.Error(ctx, "panic mid", err)
					m.Error(w, err, http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
