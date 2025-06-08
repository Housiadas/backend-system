package middleware

import (
	"errors"
	"fmt"

	"net/http"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/internal/domain/authbus"
	"github.com/Housiadas/backend-system/internal/domain/userbus"
	"github.com/Housiadas/backend-system/internal/sys/context"
	"github.com/Housiadas/backend-system/pkg/errs"
	"github.com/Housiadas/backend-system/pkg/web"
)

// UserPermissions executes authorization for resource (entity) actions
// Check if a user is allowed to modify other user's resources
func (m *Middleware) UserPermissions(rule string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userID uuid.UUID
			id := web.Param(r, "user_id")
			ctx := r.Context()

			if id != "" {
				var err error
				userID, err = uuid.Parse(id)
				if err != nil {
					err = errs.New(errs.Unauthenticated, ErrInvalidID)
					m.Log.Error(ctx, "authorize user mid: authorize", err)
					m.Error(w, err, http.StatusUnauthorized)
					return
				}

				// ensure that only one call to an expensive or duplicative operation is in flight at any given time
				response, err, _ := group.Do(fmt.Sprintf("user_id:%s", userID), func() (interface{}, error) {
					return m.Bus.User.QueryByID(ctx, userID)
				})
				if err != nil {
					switch {
					case errors.Is(err, userbus.ErrNotFound):
						err = errs.New(errs.Unauthenticated, err)
					default:
						err = errs.Newf(errs.Unauthenticated, "querybyid: userID[%s]: %s", userID, err)
					}
					m.Log.Error(ctx, "authorize user mid: authorize", err)
					m.Error(w, err, http.StatusUnauthorized)
					return
				}

				usr, ok := response.(userbus.User)
				if !ok {
					err = errs.New(errs.InternalOnlyLog, errors.New("code should be reach here"))
					m.Log.Error(ctx, "authorize error:", err)
					m.Error(w, err, http.StatusInternalServerError)
					return
				}

				// Here adds in the context the requested user based on (user_id)
				ctx = context.SetUser(ctx, usr)
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
				m.Log.Error(ctx, "authorize user mid: authorize", err)
				m.Error(w, err, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
