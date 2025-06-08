package handler

import (
	"context"
	web2 "github.com/Housiadas/backend-system/foundation/web"
	"net/http"

	"github.com/Housiadas/backend-system/app/domain/productapp"
	"github.com/Housiadas/backend-system/foundation/errs"
)

func (h *Handler) productCreate(ctx context.Context, _ http.ResponseWriter, r *http.Request) web2.Encoder {
	var app productapp.NewProduct
	if err := web2.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	prd, err := h.App.Product.Create(ctx, app)
	if err != nil {
		return errs.NewError(err)
	}

	return prd
}

func (h *Handler) productUpdate(ctx context.Context, _ http.ResponseWriter, r *http.Request) web2.Encoder {
	var app productapp.UpdateProduct
	if err := web2.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	prd, err := h.App.Product.Update(ctx, app)
	if err != nil {
		return errs.NewError(err)
	}

	return prd
}

func (h *Handler) productDelete(ctx context.Context, _ http.ResponseWriter, _ *http.Request) web2.Encoder {
	if err := h.App.Product.Delete(ctx); err != nil {
		return errs.New(errs.Internal, err)
	}

	return nil
}

func (h *Handler) productQuery(ctx context.Context, _ http.ResponseWriter, r *http.Request) web2.Encoder {
	qp := productapp.ParseQueryParams(r)

	prd, err := h.App.Product.Query(ctx, qp)
	if err != nil {
		return errs.NewError(err)
	}

	return prd
}

func (h *Handler) productQueryByID(ctx context.Context, _ http.ResponseWriter, _ *http.Request) web2.Encoder {
	prd, err := h.App.Product.QueryByID(ctx)
	if err != nil {
		return errs.NewError(err)
	}

	return prd
}
