package handlers

import (
	"context"
	"net/http"

	"github.com/Housiadas/backend-system/internal/app/service/auditapp"
	"github.com/Housiadas/backend-system/pkg/errs"
	"github.com/Housiadas/backend-system/pkg/web"
)

func (h *Handler) auditQuery(ctx context.Context, _ http.ResponseWriter, r *http.Request) web.Encoder {
	qp := auditParseQueryParams(r)

	audits, err := h.App.Audit.Query(ctx, qp)
	if err != nil {
		return errs.NewError(err)
	}

	return audits
}

func auditParseQueryParams(r *http.Request) auditapp.AppQueryParams {
	values := r.URL.Query()

	return auditapp.AppQueryParams{
		Page:      values.Get("page"),
		Rows:      values.Get("rows"),
		OrderBy:   values.Get("orderBy"),
		ObjID:     values.Get("obj_id"),
		ObjEntity: values.Get("obj_domain"),
		ObjName:   values.Get("obj_name"),
		ActorID:   values.Get("actor_id"),
		Action:    values.Get("action"),
		Since:     values.Get("since"),
		Until:     values.Get("until"),
	}
}
