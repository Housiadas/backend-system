package handlers

import (
	"context"
	"net/http"

	"github.com/Housiadas/backend-system/internal/app/usecase/transaction_usecase"
	"github.com/Housiadas/backend-system/pkg/errs"
	"github.com/Housiadas/backend-system/pkg/web"
)

func (h *Handler) transaction(ctx context.Context, _ http.ResponseWriter, r *http.Request) web.Encoder {
	var app transaction_usecase.NewTran
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	prd, err := h.App.Tx.Create(ctx, app)
	if err != nil {
		return errs.NewError(err)
	}

	return prd
}
