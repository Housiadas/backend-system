package handler

import (
	"context"
	web2 "github.com/Housiadas/backend-system/foundation/web"
	"net/http"

	"github.com/Housiadas/backend-system/app/domain/tranapp"
	"github.com/Housiadas/backend-system/foundation/errs"
)

func (h *Handler) transaction(ctx context.Context, _ http.ResponseWriter, r *http.Request) web2.Encoder {
	var app tranapp.NewTran
	if err := web2.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	prd, err := h.App.Tx.Create(ctx, app)
	if err != nil {
		return errs.NewError(err)
	}

	return prd
}
