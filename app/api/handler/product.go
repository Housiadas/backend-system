package handler

import (
	"context"
	"net/http"

	"github.com/Housiadas/backend-system/app/domain/productapp"
	"github.com/Housiadas/backend-system/business/sys/errs"
	"github.com/Housiadas/backend-system/business/web"
)

func (h *Handler) productCreate(ctx context.Context, _ http.ResponseWriter, r *http.Request) (any, error) {
	var app productapp.NewProduct
	if err := web.Decode(r, &app); err != nil {
		return nil, errs.New(errs.FailedPrecondition, err)
	}

	prd, err := h.App.Product.Create(ctx, app)
	if err != nil {
		return nil, err
	}

	return prd, nil
}

func (h *Handler) productUpdate(ctx context.Context, _ http.ResponseWriter, r *http.Request) (any, error) {
	var app productapp.UpdateProduct
	if err := web.Decode(r, &app); err != nil {
		return nil, errs.New(errs.FailedPrecondition, err)
	}

	prd, err := h.App.Product.Update(ctx, app)
	if err != nil {
		return nil, err
	}

	return prd, nil
}

func (h *Handler) productDelete(ctx context.Context, _ http.ResponseWriter, _ *http.Request) (any, error) {
	if err := h.App.Product.Delete(ctx); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) productQuery(ctx context.Context, _ http.ResponseWriter, r *http.Request) (any, error) {
	qp, err := productParseQueryParams(r)
	if err != nil {
		return nil, errs.New(errs.FailedPrecondition, err)
	}

	prd, err := h.App.Product.Query(ctx, qp)
	if err != nil {
		return nil, err
	}

	return prd, nil
}

func (h *Handler) productQueryByID(ctx context.Context, _ http.ResponseWriter, _ *http.Request) (any, error) {
	prd, err := h.App.Product.QueryByID(ctx)
	if err != nil {
		return nil, err
	}

	return prd, nil
}
