package mid

import (
	"context"
	"net/http"

	"github.com/Housiadas/backend-system/business/mid"
	"github.com/Housiadas/backend-system/foundation/logger"
	"github.com/Housiadas/backend-system/foundation/web"
)

// Errors executes the errors middleware functionality.
func Errors(log *logger.Logger) web.Middleware {
	midFunc := func(ctx context.Context, r *http.Request, next mid.Handler) (any, error) {
		return mid.Errors(ctx, log, next)
	}

	return addMiddleware(midFunc)
}
