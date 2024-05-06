package mid

import (
	"context"
	"net/http"

	"github.com/Housiadas/backend-system/business/mid"
	"github.com/Housiadas/backend-system/foundation/web"
)

// Panics executes the panic middleware functionality.
func Panics() web.Middleware {
	midFunc := func(ctx context.Context, r *http.Request, next mid.Handler) (any, error) {
		return mid.Panics(ctx, next)
	}

	return addMiddleware(midFunc)
}
