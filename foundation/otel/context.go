package otel

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type ctxKey string

const (
	tracerKey  ctxKey = "tracerKey"
	traceIDKey ctxKey = "traceIDKey"
)

func setTracer(ctx context.Context, tracer trace.Tracer) context.Context {
	return context.WithValue(ctx, tracerKey, tracer)
}

func setTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// GetTraceID returns the trace id from the context.
func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(traceIDKey).(string)
	if !ok {
		return "00000000000000000000000000000000"
	}

	return v
}
