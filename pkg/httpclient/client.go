// Package httpclient provides support for external http requests
package httpclient

import (
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/Housiadas/backend-system/pkg/logger"
)

type Config struct {
	Log     *logger.Logger
	Timeout time.Duration
}

// Client represents a http client.
type Client struct {
	log  *logger.Logger
	http *http.Client
}

// New constructs a http client
func New(cfg Config) *Client {
	cln := Client{
		log: cfg.Log,
		http: &http.Client{
			Transport: otelhttp.NewTransport(&http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 2 * time.Second, // specifies the amount of time to wait for a server's first response
			}),
			Timeout: cfg.Timeout, // specifies a time limit for requests made by this Client.
		},
	}

	return &cln
}
