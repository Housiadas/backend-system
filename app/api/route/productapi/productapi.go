// Package productapi maintains the web based api for product access.
package productapi

import (
	"context"
	"net/http"

	"github.com/Housiadas/backend-system/app/domain/productapp"
	"github.com/Housiadas/backend-system/business/sys/errs"
	"github.com/Housiadas/backend-system/foundation/web"
)

type api struct {
	productApp *productapp.Core
}

func newAPI(productApp *productapp.Core) *api {
	return &api{
		productApp: productApp,
	}
}

func (api *api) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app productapp.NewProduct
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.FailedPrecondition, err)
	}

	prd, err := api.productApp.Create(ctx, app)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, prd, http.StatusCreated)
}

func (api *api) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app productapp.UpdateProduct
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.FailedPrecondition, err)
	}

	prd, err := api.productApp.Update(ctx, app)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, prd, http.StatusOK)
}

func (api *api) delete(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	if err := api.productApp.Delete(ctx); err != nil {
		return err
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (api *api) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qp, err := parseQueryParams(r)
	if err != nil {
		return err
	}

	hme, err := api.productApp.Query(ctx, qp)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, hme, http.StatusOK)
}

func (api *api) queryByID(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	hme, err := api.productApp.QueryByID(ctx)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, hme, http.StatusOK)
}
