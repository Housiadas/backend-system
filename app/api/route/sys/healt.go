// Package sys maintains the web based api for system access.
package sys

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Housiadas/backend-system/business/data/sqldb"
	"github.com/Housiadas/backend-system/foundation/logger"
	"github.com/Housiadas/backend-system/foundation/web"
)

type api struct {
	build string
	log   *logger.Logger
	db    *sqlx.DB
}

func newApi(build string, log *logger.Logger, db *sqlx.DB) *api {
	return &api{
		build: build,
		db:    db,
		log:   log,
	}
}

// readiness checks if the database is ready and if not will return a 500 status.
// Do not respond by just returning an error because further up in the call
// stack it will interpret that as a non-trusted error.
func (a *api) readiness(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	status := "ok"
	statusCode := http.StatusOK
	err := sqldb.StatusCheck(ctx, a.db)
	if err != nil {
		status = "db not ready"
		statusCode = http.StatusInternalServerError
		a.log.Info(ctx, "readiness failure", "status", status)
	}

	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	return web.Respond(ctx, w, data, statusCode)
}

// liveness returns simple status info if the service is alive. If the
// app is deployed to a Kubernetes cluster, it will also return pod, node, and
// namespace details via the Downward API. The Kubernetes environment variables
// need to be set within your Pod/Deployment manifest.
func (a *api) liveness(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	data := struct {
		Status     string `json:"status,omitempty"`
		Build      string `json:"build,omitempty"`
		Host       string `json:"host,omitempty"`
		Name       string `json:"name,omitempty"`
		PodIP      string `json:"podIP,omitempty"`
		Node       string `json:"node,omitempty"`
		Namespace  string `json:"namespace,omitempty"`
		GOMAXPROCS int    `json:"GOMAXPROCS,omitempty"`
	}{
		Status:     "up",
		Build:      a.build,
		Host:       host,
		Name:       os.Getenv("KUBERNETES_NAME"),
		PodIP:      os.Getenv("KUBERNETES_POD_IP"),
		Node:       os.Getenv("KUBERNETES_NODE_NAME"),
		Namespace:  os.Getenv("KUBERNETES_NAMESPACE"),
		GOMAXPROCS: runtime.GOMAXPROCS(0),
	}

	// This handler provides a free timer loop.

	return web.Respond(ctx, w, data, http.StatusOK)
}
