package handler

import (
	"context"
	"net/http"

	"github.com/Housiadas/backend-system/app/domain/tranapp"
	"github.com/Housiadas/backend-system/business/sys/errs"
	"github.com/Housiadas/backend-system/business/web"
)

func (h *Handler) transaction(ctx context.Context, _ http.ResponseWriter, r *http.Request) web.Encoder {
	var app tranapp.NewTran
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	prd, err := h.App.Tx.Create(ctx, app)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return prd
}
