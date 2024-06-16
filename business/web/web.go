package web

import (
	"context"
	"net/http"
)

// HandlerFunc represents a function that handles a http request
type HandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) (Encoder, error)

// Encoder defines behavior that can encode a data model and provide
// the content type for that encoding.
type Encoder interface {
	Encode() (data []byte, contentType string, err error)
}

type httpStatus interface {
	HTTPStatus() int
}
