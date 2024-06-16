package mid

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/business/sys/errs"
	"github.com/Housiadas/backend-system/business/web"
)

// Bearer processes JWT token.
func (m *Mid) Bearer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			claims, err := m.Auth.Authenticate(ctx, r.Header.Get("authorization"))
			if err != nil {
				m.Log.Error(ctx, "bearer mid: unauthenticated", errs.Unauthenticated)
				return
			}

			if claims.Subject == "" {
				m.Log.Info(ctx, "request unauthenticated", errs.Unauthenticated)
				http.Error(w, errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, no claims").Error(), errs.Unauthenticated.Value())
				return
			}

			subjectID, err := uuid.Parse(claims.Subject)
			if err != nil {
				m.Log.Error(ctx, "bearer mid: parsing", errs.Newf(errs.Unauthenticated, "parsing subject: %s", err))
				return
			}

			ctx = web.SetUserID(ctx, subjectID)
			ctx = web.SetClaims(ctx, claims)

			next.ServeHTTP(w, r)
		})
	}
}
