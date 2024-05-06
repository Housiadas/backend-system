package mid

import (
	"context"

	"github.com/Housiadas/backend-system/foundation/metrics"
)

// Metrics updates program counters.
func Metrics(ctx context.Context, next Handler) (any, error) {
	ctx = metrics.Set(ctx)

	resp, err := next(ctx)

	n := metrics.AddRequests(ctx)

	if n%1000 == 0 {
		metrics.AddGoroutines(ctx)
	}

	if err != nil {
		metrics.AddErrors(ctx)
	}

	return resp, err
}
