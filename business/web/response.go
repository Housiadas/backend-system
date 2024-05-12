package web

import (
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel/attribute"

	"github.com/Housiadas/backend-system/foundation/logger"
	"github.com/Housiadas/backend-system/foundation/tracer"
)

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

func (resp *Respond) Respond(data any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		statusCode := http.StatusOK

		switch v := data.(type) {
		case httpStatus:
			statusCode = v.HTTPStatus()
		case error:
			statusCode = http.StatusInternalServerError
		}

		_, span := tracer.AddSpan(ctx, "web.response", attribute.Int("status", statusCode))
		defer span.End()

		if data == nil {
			statusCode = http.StatusNoContent
		}

		if statusCode == http.StatusNoContent {
			w.WriteHeader(statusCode)
			return
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			resp.Log.Error(ctx, "web.respond: marshal error", err)
			return
		}

		w.WriteHeader(statusCode)
		if _, err := w.Write(jsonData); err != nil {
			resp.Log.Error(ctx, "web.respond: write error", err)
			return
		}

		return
	}
}
