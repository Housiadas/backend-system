package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"

	"go.opentelemetry.io/otel/attribute"

	"github.com/Housiadas/backend-system/foundation/errs"
	"github.com/Housiadas/backend-system/foundation/logger"
	"github.com/Housiadas/backend-system/foundation/otel"
)

// HandlerFunc represents a function that handles an http request
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

		// Get status code
		statusCode := respond.statusCode(resp)

		// Record errors with status code 500 and above
		err := isError(resp)
		if err != nil {
			resp = respond.errorRecorder(ctx, statusCode, err)
		}

		// Send response back to a client
		if err := respond.response(ctx, w, statusCode, resp); err != nil {
			respond.Log.Error(ctx, "web-respond", "ERROR", err)
		}
	}
}

func (respond *Respond) response(ctx context.Context, w http.ResponseWriter, statusCode int, dataModel Encoder) error {
	// If the context has been canceled, it means the client is no longer waiting for a response.
	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			return errors.New("client disconnected, do not send response")
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

func (respond *Respond) errorRecorder(ctx context.Context, statusCode int, err error) Encoder {
	var appErr *errs.Error
	ok := errors.As(err, &appErr)
	if !ok {
		appErr = errs.Newf(errs.Internal, "Internal Server Error")
	}

	// If not, critical error does not record it
	if statusCode < 500 {
		return appErr
	}

	_, span := otel.AddSpan(ctx, "app.response.error")
	span.RecordError(err)
	defer span.End()

	respond.Log.Error(ctx, "error during request",
		"err", err,
		"source_err_file", path.Base(appErr.FileName),
		"source_err_func", path.Base(appErr.FuncName),
	)

	if appErr.Code == errs.InternalOnlyLog {
		appErr = errs.Newf(errs.Internal, "Internal Server Error")
	}

	// Send the error back so it can be used as the response.
	return appErr
}

// isError checks if the Encoder has an error inside it.
func isError(e Encoder) error {
	err, isError := e.(error)
	if isError {
		return err
	}
	return nil
}

func (respond *Respond) statusCode(dataModel Encoder) int {
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

	return statusCode
}
