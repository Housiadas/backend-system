package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/internal/domain/authbus"
	"github.com/Housiadas/backend-system/internal/domain/productbus"
	"github.com/Housiadas/backend-system/internal/sys/context"
	"github.com/Housiadas/backend-system/pkg/errs"
	"github.com/Housiadas/backend-system/pkg/web"
)

// ProductPermissions executes authorization for resource (entity) actions
// Check if a user is allowed to modify product
func (m *Middleware) ProductPermissions(rule string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userID uuid.UUID
			id := web.Param(r, "product_id")
			ctx := r.Context()

			if id != "" {
				var err error
				productID, err := uuid.Parse(id)
				if err != nil {
					err = errs.New(errs.Unauthenticated, ErrInvalidID)
					m.Log.Error(ctx, "authorize product mid: authorize", err)
					m.Error(w, err, http.StatusUnauthorized)
					return
				}

				// ensure that only one call to an expensive or duplicative operation is in flight at any given time
				response, err, _ := group.Do(fmt.Sprintf("product_id:%s", productID), func() (interface{}, error) {
					return m.Bus.Product.QueryByID(ctx, productID)
				})
				if err != nil {
					switch {
					case errors.Is(err, productbus.ErrNotFound):
						err = errs.New(errs.Unauthenticated, err)
					default:
						err = errs.Newf(errs.Internal, "querybyid: productID[%s]: %s", productID, err)
					}
					m.Log.Error(ctx, "authorize product mid: authorize", err)
					m.Error(w, err, http.StatusUnauthorized)
					return
				}

				prd, ok := response.(productbus.Product)
				if !ok {
					err = errs.New(errs.InternalOnlyLog, errors.New("code should be reach here"))
					m.Log.Error(ctx, "authorize error:", err)
					m.Error(w, err, http.StatusInternalServerError)
				}

				userID = prd.UserID
				ctx = context.SetProduct(ctx, prd)
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
				m.Log.Error(ctx, "authorize product mid: authorize", err)
				m.Error(w, err, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
