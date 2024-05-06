// Package systemapi maintains the web based api for system access.
package systemapi

import (
	"context"
	"net/http"

	"github.com/Housiadas/backend-system/app/domain/systemapp"
	"github.com/Housiadas/backend-system/business/sys/errs"
)

type api struct {
	systemApp *systemapp.App
}

func newAPI(systemApp *systemapp.App) *api {
	return &api{
		systemApp: systemApp,
	}
}

// readiness checks if the database is ready and if not will return a 500 status.
// Do not respond by just returning an error because further up in the call
// stack it will interpret that as a non-trusted error.
func (api *api) readiness(ctx context.Context, _ http.ResponseWriter, _ *http.Request) (any, error) {
	if err := api.systemApp.Readiness(ctx); err != nil {
		return nil, errs.Newf(errs.Internal, "database not ready")
	}

	data := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	return data, nil
}

// liveness returns simple status info if the service is alive. If the
// app is deployed to a Kubernetes cluster, it will also return pod, node, and
// namespace details via the Downward API. The Kubernetes environment variables
// need to be set within your Pod/Deployment manifest.
func (api *api) liveness(_ context.Context, _ http.ResponseWriter, _ *http.Request) (any, error) {
	info := api.systemApp.Liveness()

	return info, nil
}
