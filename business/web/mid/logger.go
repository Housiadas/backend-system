package mid

import (
	"fmt"
	"net/http"
	"time"
)

// Logger writes information about the request to the logs.
func (m *Mid) Logger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ctx := r.Context()
			path := r.URL.Path
			method := r.Method
			remoteAddr := r.RemoteAddr

			rawQuery := r.URL.RawQuery
			if rawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, rawQuery)
			}

			m.Log.Info(ctx, "request started", "method", method, "path", path, "remoteaddr", remoteAddr)

			defer func() {
				m.Log.Info(ctx, "request completed",
					"path", path,
					"method", method,
					"remote_addr", remoteAddr,
					"user_agent", r.UserAgent(),
					"execution_time", time.Since(start).String(),
				)
			}()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
