package middleware

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/business/sys/context"
	"github.com/Housiadas/backend-system/pkg/errs"
)

// Bearer processes JWT token.
func (m *Middleware) Bearer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			claims, err := m.Bus.Auth.Authenticate(ctx, r.Header.Get("authorization"))
			if err != nil {
				err = errs.New(errs.Unauthenticated, err)
				m.Log.Error(ctx, "bearer mid: unauthenticated", errs.Unauthenticated)
				m.Error(w, err, http.StatusUnauthorized)
				return
			}

			if claims.Subject == "" {
				err = errs.New(errs.Unauthenticated,
					errors.New("authorize: you are not authorized for that action, no claims"),
				)
				m.Log.Info(ctx, "request unauthenticated", err)
				m.Error(w, err, http.StatusUnauthorized)
				return
			}

			subjectID, err := uuid.Parse(claims.Subject)
			if err != nil {
				m.Log.Error(ctx,
					"bearer mid: parsing",
					errs.Newf(errs.Unauthenticated, "parsing subject: %s", err),
				)
				return
			}

			ctx = context.SetClaims(ctx, claims)
			ctx = context.SetUserID(ctx, subjectID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
