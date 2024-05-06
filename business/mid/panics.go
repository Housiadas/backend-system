package mid

import (
	"context"
	"runtime/debug"

	"github.com/Housiadas/backend-system/business/sys/errs"
	"github.com/Housiadas/backend-system/foundation/metrics"
)

// Panics recovers from panics and converts the panic to an error, so it is
// reported in Metrics and handled in Errors.
func Panics(ctx context.Context, next Handler) (resp any, err error) {

	// Defer a function to recover from a panic and set the err return
	// variable after the fact.
	defer func() {
		if rec := recover(); rec != nil {
			trace := debug.Stack()
			err = errs.Newf(errs.Internal, "PANIC [%v] TRACE[%s]", rec, string(trace))

			metrics.AddPanics(ctx)
		}
	}()

	return next(ctx)
}
