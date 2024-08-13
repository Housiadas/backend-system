package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"

	"github.com/Housiadas/backend-system/foundation/logger"
	"github.com/Housiadas/backend-system/foundation/otel"
)

// HandlerFunc represents a function that handles a http request
type HandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) Encoder

// Encoder defines behavior that can encode a data model and provide the content type for that encoding.
type Encoder interface {
	Encode() (data []byte, contentType string, err error)
}

type httpStatus interface {
	HTTPStatus() int
}

type Respond struct {
	Log *logger.Logger
}

func NewRespond(log *logger.Logger) *Respond {
	return &Respond{
		Log: log,
	}
}

func (respond *Respond) Respond(handlerFunc HandlerFunc) http.HandlerFunc {
	// This is the decorator/middleware pattern in golang
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Executes the handlerFunc for the specific route
		resp := handlerFunc(ctx, w, r)
		if err := response(ctx, w, resp); err != nil {
			respond.Log.Error(ctx, "web-respond", "ERROR", err)
		}
	}
}

func response(ctx context.Context, w http.ResponseWriter, dataModel Encoder) error {
	// If the context has been canceled, it means the client is no longer waiting for a response.
	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			return errors.New("client disconnected, do not send response")
		}
	}

	var statusCode = http.StatusOK

	switch v := dataModel.(type) {
	case httpStatus:
		statusCode = v.HTTPStatus()
	case error:
		statusCode = http.StatusInternalServerError
	default:
		if dataModel == nil {
			statusCode = http.StatusNoContent
		}
	}

	_, span := otel.AddSpan(ctx, "web.response", attribute.Int("status", statusCode))
	defer span.End()

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	data, contentType, err := dataModel.Encode()
	if err != nil {
		return fmt.Errorf("respond: encode: %w", err)
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("respond: write: %w", err)
	}

	return nil
}
