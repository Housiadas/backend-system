package mid

import (
	"net/http"

	"github.com/Housiadas/backend-system/foundation/otel"
)

// Otel starts the otel tracing and stores the trace id in the context.
func (m *Mid) Otel() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = otel.InjectTracing(ctx, m.Tracer)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
