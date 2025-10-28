package handlers

import (
	"context"
	"net/http"

	"github.com/Housiadas/backend-system/internal/app/usecase/product_usecase"
	"github.com/Housiadas/backend-system/pkg/errs"
	"github.com/Housiadas/backend-system/pkg/web"
)

func (h *Handler) productCreate(ctx context.Context, _ http.ResponseWriter, r *http.Request) web.Encoder {
	var app product_usecase.NewProduct
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	prd, err := h.App.Product.Create(ctx, app)
	if err != nil {
		return errs.NewError(err)
	}

	return prd
}

func (h *Handler) productUpdate(ctx context.Context, _ http.ResponseWriter, r *http.Request) web.Encoder {
	var app product_usecase.UpdateProduct
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	prd, err := h.App.Product.Update(ctx, app)
	if err != nil {
		return errs.NewError(err)
	}

	return prd
}

func (h *Handler) productDelete(ctx context.Context, _ http.ResponseWriter, _ *http.Request) web.Encoder {
	if err := h.App.Product.Delete(ctx); err != nil {
		return errs.New(errs.Internal, err)
	}

	return nil
}

func (h *Handler) productQuery(ctx context.Context, _ http.ResponseWriter, r *http.Request) web.Encoder {
	qp := productParseQueryParams(r)

	prd, err := h.App.Product.Query(ctx, qp)
	if err != nil {
		return errs.NewError(err)
	}

	return prd
}

func (h *Handler) productQueryByID(ctx context.Context, _ http.ResponseWriter, _ *http.Request) web.Encoder {
	prd, err := h.App.Product.QueryByID(ctx)
	if err != nil {
		return errs.NewError(err)
	}

	return prd
}

func productParseQueryParams(r *http.Request) product_usecase.AppQueryParams {
	values := r.URL.Query()

	return product_usecase.AppQueryParams{
		Page:     values.Get("page"),
		Rows:     values.Get("rows"),
		OrderBy:  values.Get("orderBy"),
		ID:       values.Get("product_id"),
		Name:     values.Get("name"),
		Cost:     values.Get("cost"),
		Quantity: values.Get("quantity"),
	}
}
