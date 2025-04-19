// Package mid provides app level mid support.
package mid

import (
	"bytes"
	"encoding/json"
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

type Config struct {
	Log     *logger.Logger
	Tracer  trace.Tracer
	Tx      *sqldb.DBBeginner
	Auth    *authbus.Auth
	User    *userbus.Business
	Product *productbus.Business
}

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

func New(cfg Config) *Mid {
	return &Mid{
		Bus: Business{
			Auth:    cfg.Auth,
			User:    cfg.User,
			Product: cfg.Product,
		},
		Log:    cfg.Log,
		Tracer: cfg.Tracer,
		Tx:     cfg.Tx,
	}
}

func (m *Mid) Error(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(err); err != nil {
		return
	}
	return
}

// ResponseRecorder a custom http.ResponseWriter to capture the response before it's sent to the client.
// We are capturing the result of the handlers to the middleware
type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rec *ResponseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// Capture the response body
func (rec *ResponseRecorder) Write(b []byte) (int, error) {
	rec.body.Write(b)
	return rec.ResponseWriter.Write(b)
}
