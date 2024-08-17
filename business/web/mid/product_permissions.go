package mid

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/business/domain/authbus"
	"github.com/Housiadas/backend-system/business/domain/productbus"
	"github.com/Housiadas/backend-system/business/sys/errs"
	"github.com/Housiadas/backend-system/business/web"
)

// ProductPermissions executes authorization for resource (entity) actions
// Check if a user is allowed to modify product
func (m *Mid) ProductPermissions(rule string) func(next http.Handler) http.Handler {
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

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					if err := json.NewEncoder(w).Encode(err); err != nil {
						return
					}
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

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					if err := json.NewEncoder(w).Encode(err); err != nil {
						return
					}
					return
				}

				prd, ok := response.(productbus.Product)
				if !ok {
					err = errs.New(errs.InternalOnlyLog, errors.New("code should be reach here"))
					m.Log.Error(ctx, "authorize error:", err)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					if err := json.NewEncoder(w).Encode(err); err != nil {
						return
					}
					return
				}

				userID = prd.UserID
				ctx = web.SetProduct(ctx, prd)
			}

			authData := authbus.Authorize{
				Claims: web.GetClaims(ctx),
				UserID: userID,
				Rule:   rule,
			}

			if err := m.Bus.Auth.Authorize(ctx, authData.Claims, authData.UserID, authData.Rule); err != nil {
				err = errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", authData.Claims.Roles, authData.Rule, err)
				m.Log.Error(ctx, "authorize product mid: authorize", err)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				if err := json.NewEncoder(w).Encode(err); err != nil {
					return
				}
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
