// Package systemapi maintains the web based api for system access.
package systemapi

import (
	"context"
	"github.com/Housiadas/backend-system/app/domain/systemapp"
	"github.com/Housiadas/backend-system/business/sys/errs"
	"net/http"
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
//
// Readiness godoc
// @Summary      App Readiness
// @Description  Check application's readiness
// @Accept       json
// @Produce      json
// @Success      200  {object}  systemapp.Status
// @Failure      500  {object}  errs.Error
// @Router       /readiness [get]
func (api *api) readiness(ctx context.Context, _ http.ResponseWriter, _ *http.Request) (any, error) {
	if err := api.systemApp.Readiness(ctx); err != nil {
		return nil, errs.Newf(errs.Internal, "database not ready")
	}

	data := systemapp.Status{
		Status: "OK",
	}

	return data, nil
}

// liveness returns simple status info if the service is alive. If the
// app is deployed to a Kubernetes cluster, it will also return pod, node, and
// namespace details via the Downward API. The Kubernetes environment variables
// need to be set within your Pod/Deployment manifest.
//
// Liveness godoc
// @Summary      App Liveness
// @Description  Returns application's status info if the service is alive
// @Accept       json
// @Produce      json
// @Success      200  {object}  systemapp.Info
// @Router       /liveness [get]
func (api *api) liveness(_ context.Context, _ http.ResponseWriter, _ *http.Request) (any, error) {
	info := api.systemApp.Liveness()

	return info, nil
}

// Swagger http handler
func (api *api) swagger(ctx context.Context, w http.ResponseWriter, r *http.Request) (any, error) {
	return api.systemApp.Swagger(ctx, w, r)
}
