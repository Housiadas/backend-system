package mid

import (
	"context"
	"net/http"

	"github.com/Housiadas/backend-system/business/auth"
	"github.com/Housiadas/backend-system/business/mid"
	"github.com/Housiadas/backend-system/foundation/web"
)

// Authenticate processes JWT authentication logic.
func Authenticate(ath *auth.Auth) web.Middleware {
	midFunc := func(ctx context.Context, r *http.Request, next mid.Handler) (any, error) {
		return mid.Bearer(ctx, ath, r.Header.Get("authorization"), next)
	}

	return addMiddleware(midFunc)
}
