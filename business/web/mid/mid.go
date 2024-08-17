// Package mid provides app level mid support.
package mid

import (
	"bytes"
	"errors"
	"net/http"

	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/singleflight"

	"github.com/Housiadas/backend-system/business/data/sqldb"
	"github.com/Housiadas/backend-system/business/domain/authbus"
	"github.com/Housiadas/backend-system/business/domain/productbus"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/foundation/logger"
)

var (
	// ErrInvalidID represents a condition where the id is not an uuid.
	ErrInvalidID = errors.New("ID is not in its proper form")

	group = singleflight.Group{}
)

type Mid struct {
	Bus    Business
	Log    *logger.Logger
	Tracer trace.Tracer
	Tx     *sqldb.DBBeginner
}

type Business struct {
	Auth    *authbus.Auth
	User    *userbus.Business
	Product *productbus.Business
}

func New(b Business, l *logger.Logger, t trace.Tracer, tx *sqldb.DBBeginner) *Mid {
	return &Mid{
		Bus:    b,
		Log:    l,
		Tracer: t,
		Tx:     tx,
	}
}

// ResponseRecorder a custom http.ResponseWriter to capture the response
// before it's sent to the client. We are capturing the result of the handlers to the middleware
type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rec *ResponseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *ResponseRecorder) Write(b []byte) (int, error) {
	rec.body.Write(b) // Capture the response body
	return rec.ResponseWriter.Write(b)
}
