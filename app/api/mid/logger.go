package mid

import (
	"context"
	"net/http"

	"github.com/Housiadas/backend-system/business/mid"
	"github.com/Housiadas/backend-system/foundation/logger"
	"github.com/Housiadas/backend-system/foundation/web"
)

// Logger executes the logger middleware functionality.
func Logger(log *logger.Logger) web.Middleware {
	midFunc := func(ctx context.Context, r *http.Request, next mid.Handler) (any, error) {
		return mid.Logger(ctx, log, r.URL.Path, r.URL.RawQuery, r.Method, r.RemoteAddr, next)
	}

	return addMiddleware(midFunc)
}
