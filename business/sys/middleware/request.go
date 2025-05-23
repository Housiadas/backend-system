package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/business/sys/web"
)

// RequestID is a middleware that injects uuid as middleware.RequestIDHeader when not present
func (m *Middleware) RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var u uuid.UUID
		var err error

		ctx := r.Context()
		h := r.Header.Get(middleware.RequestIDHeader)

		if h == "" {
			u = uuid.New()
		} else {
			u, err = uuid.Parse(h)
			if err != nil {
				m.Log.Info(ctx, "request id parse error", err)
				u = uuid.New()
			}
		}

		us := u.String()
		ctx = web.SetRequestID(ctx, us)
		w.Header().Set(middleware.RequestIDHeader, us)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
