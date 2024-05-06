package mid

import (
	"context"
	"net/http"

	"github.com/Housiadas/backend-system/business/auth"
	"github.com/Housiadas/backend-system/business/domain/productbus"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/mid"
	"github.com/Housiadas/backend-system/foundation/web"
)

// Authorize validates authorization via the auth service.
func Authorize(ath *auth.Auth, rule string) web.Middleware {
	midFunc := func(ctx context.Context, r *http.Request, next mid.Handler) (any, error) {
		return mid.Authorize(ctx, ath, rule, next)
	}

	return addMiddleware(midFunc)
}

// AuthorizeUser executes the specified role and extracts the specified
// user from the DB if a user id is specified in the call. Depending on the rule
// specified, the userid from the claims may be compared with the specified
// user id.
func AuthorizeUser(ath *auth.Auth, userBus *userbus.Business, rule string) web.Middleware {
	midFunc := func(ctx context.Context, r *http.Request, next mid.Handler) (any, error) {
		return mid.AuthorizeUser(ctx, ath, userBus, rule, web.Param(r, "user_id"), next)
	}

	return addMiddleware(midFunc)
}

// AuthorizeProduct executes the specified role and extracts the specified
// product from the DB if a product id is specified in the call. Depending on
// the rule specified, the userid from the claims may be compared with the
// specified user id from the product.
func AuthorizeProduct(ath *auth.Auth, productBus *productbus.Business) web.Middleware {
	midFunc := func(ctx context.Context, r *http.Request, next mid.Handler) (any, error) {
		return mid.AuthorizeProduct(ctx, ath, productBus, web.Param(r, "product_id"), next)
	}

	return addMiddleware(midFunc)
}
