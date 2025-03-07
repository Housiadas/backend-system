package handler

import (
	"context"
	"net/http"

	"github.com/Housiadas/backend-system/app/domain/systemapp"
	"github.com/Housiadas/backend-system/business/sys/errs"
	"github.com/Housiadas/backend-system/business/web"
)

// readiness checks if the database is ready and if not will return a 500 status.
// Do not respond by just returning an error because further up in the call
// stack it will interpret that as a non-trusted error.
//
// Readiness godoc
// @Summary      App Readiness
// @Description  Check application's readiness
// @Tags		 System
// @Accept       json
// @Produce      json
// @Success      200  {object}  systemapp.Status
// @Failure      500  {object}  errs.Error
// @Router       /readiness [get]
func (h *Handler) readiness(ctx context.Context, _ http.ResponseWriter, _ *http.Request) web.Encoder {
	if err := h.App.System.Readiness(ctx); err != nil {
		return errs.Newf(errs.Internal, "database not ready")
	}

	data := systemapp.Status{
		Status: "OK",
	}

	return data
}

// liveness returns simple status info if the service is alive. If the
// app is deployed to a Kubernetes cluster, it will also return pod, node, and
// namespace details via the Downward API. The Kubernetes environment variables
// need to be set within your Pod/Deployment manifest.
//
// Liveness godoc
// @Summary      App Liveness
// @Description  Returns application's status info if the service is alive
// @Tags		 System
// @Accept       json
// @Produce      json
// @Success      200  {object}  systemapp.Info
// @Router       /liveness [get]
func (h *Handler) liveness(_ context.Context, _ http.ResponseWriter, _ *http.Request) web.Encoder {
	info := h.App.System.Liveness()

	return info
}
