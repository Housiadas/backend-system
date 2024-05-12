package mid

import (
	"net/http"

	"github.com/Housiadas/backend-system/business/web"
)

func (m *Mid) ApiVersion(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := web.SetApiVersion(r.Context(), version)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
