package mid

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/business/auth"
	"github.com/Housiadas/backend-system/business/domain/productbus"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/sys/errs"
)

// ErrInvalidID represents a condition where the id is not a uuid.
var ErrInvalidID = errors.New("ID is not in its proper form")

// Authorize validates authorization.
func Authorize(ctx context.Context, ath *auth.Auth, rule string, next Handler) (any, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, errs.New(errs.Unauthenticated, err)
	}

	authData := auth.Authorize{
		Claims: GetClaims(ctx),
		UserID: userID,
		Rule:   rule,
	}

	if err := ath.Authorize(ctx, authData.Claims, authData.UserID, authData.Rule); err != nil {
		return nil, errs.New(errs.Unauthenticated, err)
	}

	return next(ctx)
}

// AuthorizeUser executes the specified role and extracts the specified
// user from the DB if a user id is specified in the call. Depending on the rule
// specified, the userid from the claims may be compared with the specified
// user id.
func AuthorizeUser(ctx context.Context, ath *auth.Auth, userBus *userbus.Business, rule string, id string, next Handler) (any, error) {
	var userID uuid.UUID

	if id != "" {
		var err error
		userID, err = uuid.Parse(id)
		if err != nil {
			return nil, errs.New(errs.Unauthenticated, ErrInvalidID)
		}

		usr, err := userBus.QueryByID(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, userbus.ErrNotFound):
				return nil, errs.New(errs.Unauthenticated, err)
			default:
				return nil, errs.Newf(errs.Unauthenticated, "querybyid: userID[%s]: %s", userID, err)
			}
		}

		ctx = setUser(ctx, usr)
	}

	authData := auth.Authorize{
		Claims: GetClaims(ctx),
		UserID: userID,
		Rule:   rule,
	}

	if err := ath.Authorize(ctx, authData.Claims, authData.UserID, authData.Rule); err != nil {
		return nil, errs.New(errs.Unauthenticated, err)
	}

	return next(ctx)
}

// AuthorizeProduct executes the specified role and extracts the specified
// product from the DB if a product id is specified in the call. Depending on
// the rule specified, the userid from the claims may be compared with the
// specified user id from the product.
func AuthorizeProduct(ctx context.Context, ath *auth.Auth, productBus *productbus.Business, id string, next Handler) (any, error) {
	var userID uuid.UUID

	if id != "" {
		var err error
		productID, err := uuid.Parse(id)
		if err != nil {
			return nil, errs.New(errs.Unauthenticated, ErrInvalidID)
		}

		prd, err := productBus.QueryByID(ctx, productID)
		if err != nil {
			switch {
			case errors.Is(err, productbus.ErrNotFound):
				return nil, errs.New(errs.Unauthenticated, err)
			default:
				return nil, errs.Newf(errs.Internal, "querybyid: productID[%s]: %s", productID, err)
			}
		}

		userID = prd.UserID
		ctx = setProduct(ctx, prd)
	}

	authData := auth.Authorize{
		UserID: userID,
		Claims: GetClaims(ctx),
		Rule:   auth.RuleAdminOrSubject,
	}

	if err := ath.Authorize(ctx, authData.Claims, authData.UserID, authData.Rule); err != nil {
		return nil, errs.New(errs.Unauthenticated, err)
	}

	return next(ctx)
}
