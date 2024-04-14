package mid

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/business/auth"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/mid"
	"github.com/Housiadas/backend-system/business/sys/errs"
	"github.com/Housiadas/backend-system/foundation/web"
)

// ErrInvalidID represents a condition where the id is not a uuid.
var ErrInvalidID = errors.New("ID is not in its proper form")

// Authorize executes the specified role and does not extract any domain data.
func Authorize(a *auth.Auth, rule string) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			userID, err := mid.GetUserID(ctx)
			if err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			authData := auth.Authorize{
				Claims: mid.GetClaims(ctx),
				UserID: userID,
				Rule:   rule,
			}

			if err := a.Authorize(ctx, authData.Claims, authData.UserID, authData.Rule); err != nil {
				return errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", authData.Claims.Roles, authData.Rule, err)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

// AuthorizeUser executes the specified role and extracts the specified user
// from the DB if a user id is specified in the call. Depending on the rule
// specified, the userid from the claims may be compared with the specified
// user id.
func AuthorizeUser(a *auth.Auth, userBus *userbus.Core, rule string) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var userID uuid.UUID

			if id := web.Param(r, "user_id"); id != "" {
				var err error
				userID, err = uuid.Parse(id)
				if err != nil {
					return errs.New(errs.Unauthenticated, ErrInvalidID)
				}

				usr, err := userBus.QueryByID(ctx, userID)
				if err != nil {
					switch {
					case errors.Is(err, userbus.ErrNotFound):
						return errs.New(errs.Unauthenticated, err)
					default:
						return errs.Newf(errs.Unauthenticated, "querybyid: userID[%s]: %s", userID, err)
					}
				}

				ctx = mid.SetUser(ctx, usr)
			}

			authData := auth.Authorize{
				Claims: mid.GetClaims(ctx),
				UserID: userID,
				Rule:   rule,
			}

			if err := a.Authorize(ctx, authData.Claims, authData.UserID, authData.Rule); err != nil {
				return errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", authData.Claims.Roles, authData.Rule, err)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
