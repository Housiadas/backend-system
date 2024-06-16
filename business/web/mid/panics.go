package mid

import (
	"net/http"
	"runtime/debug"

	"github.com/Housiadas/backend-system/business/sys/errs"
	"github.com/Housiadas/backend-system/foundation/metrics"
)

// Recoverer recovers from panics and converts the panic to an error, so it is reported in Metrics and handled in Errors.
func (m *Mid) Recoverer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Defer a function to recover from panic
			defer func() {
				if rec := recover(); rec != nil {
					ctx := r.Context()
					metrics.AddPanics(ctx)
					trace := debug.Stack()
					err := errs.Newf(errs.Internal, "PANIC [%v] TRACE[%s]", rec, string(trace))
					m.Log.Error(ctx, "panic mid", err)
					http.Error(w, err.Error(), errs.Internal.Value())
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
