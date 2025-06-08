package middleware

import (
	"net/http"

	"github.com/Housiadas/backend-system/business/domain/authbus"
	"github.com/Housiadas/backend-system/business/sys/context"
	"github.com/Housiadas/backend-system/foundation/errs"
)

// Authorize validates user's role.
func (m *Middleware) Authorize(rule string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userID, err := context.GetUserID(ctx)
			if err != nil {
				err = errs.New(errs.Unauthenticated, err)
				m.Log.Error(ctx, "authorize mid: get user id", err)
				m.Error(w, err, http.StatusUnauthorized)
				return
			}

			authData := authbus.Authorize{
				Claims: context.GetClaims(ctx),
				UserID: userID,
				Rule:   rule,
			}

			if err := m.Bus.Auth.Authorize(ctx, authData.Claims, authData.UserID, authData.Rule); err != nil {
				err = errs.Newf(errs.Unauthenticated,
					"authorize: you are not authorized for that action, claims[%v] rule[%v]: %s",
					authData.Claims.Roles, authData.Rule, err,
				)
				m.Log.Error(ctx, "authorize mid: authorize", err)
				m.Error(w, err, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
