package mid

import (
	"context"
	"net/http"

	"github.com/Housiadas/backend-system/business/mid"
	"github.com/Housiadas/backend-system/foundation/web"
)

// Metrics updates program counters using the middleware functionality.
func Metrics() web.Middleware {
	midFunc := func(ctx context.Context, r *http.Request, next mid.Handler) (any, error) {
		return mid.Metrics(ctx, next)
	}

	return addMiddleware(midFunc)
}
