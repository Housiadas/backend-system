// Package product_usecase maintains the cli layer api for the product core.
package product_usecase

import (
	"context"

	ctxPck "github.com/Housiadas/backend-system/internal/common/context"
	"github.com/Housiadas/backend-system/internal/common/validation"
	"github.com/Housiadas/backend-system/internal/core/service/productcore"
	"github.com/Housiadas/backend-system/pkg/errs"
	"github.com/Housiadas/backend-system/pkg/order"
	"github.com/Housiadas/backend-system/pkg/page"
)

// App manages the set of cli layer api functions for the product core.
type App struct {
	productBus *productcore.Core
}

// NewApp constructs a product cli API for use.
func NewApp(productBus *productcore.Core) *App {
	return &App{
		productBus: productBus,
	}
}

// Create adds a new product to the system.
func (a *App) Create(ctx context.Context, app NewProduct) (Product, error) {
	np, err := toBusNewProduct(ctx, app)
	if err != nil {
		return Product{}, errs.New(errs.InvalidArgument, err)
	}

	prd, err := a.productBus.Create(ctx, np)
	if err != nil {
		return Product{}, errs.Newf(errs.Internal, "create: prd[%+v]: %s", prd, err)
	}

	return toAppProduct(prd), nil
}

// Update updates an existing product.
func (a *App) Update(ctx context.Context, app UpdateProduct) (Product, error) {
	up, err := toBusUpdateProduct(app)
	if err != nil {
		return Product{}, errs.New(errs.InvalidArgument, err)
	}

	prd, err := ctxPck.GetProduct(ctx)
	if err != nil {
		return Product{}, errs.Newf(errs.Internal, "product missing in context: %s", err)
	}

	updPrd, err := a.productBus.Update(ctx, prd, up)
	if err != nil {
		return Product{}, errs.Newf(errs.Internal, "update: productID[%s] up[%+v]: %s", prd.ID, app, err)
	}

	return toAppProduct(updPrd), nil
}

// Delete removes a product from the system.
func (a *App) Delete(ctx context.Context) error {
	prd, err := ctxPck.GetProduct(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "productID missing in context: %s", err)
	}

	if err := a.productBus.Delete(ctx, prd); err != nil {
		return errs.Newf(errs.Internal, "delete: productID[%s]: %s", prd.ID, err)
	}

	return nil
}

// Query returns a list of products with paging.
func (a *App) Query(ctx context.Context, qp AppQueryParams) (page.Result[Product], error) {
	p, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return page.Result[Product]{}, validation.NewFieldErrors("page", err)
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return page.Result[Product]{}, err
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, defaultOrderBy)
	if err != nil {
		return page.Result[Product]{}, validation.NewFieldErrors("order", err)
	}

	products, err := a.productBus.Query(ctx, filter, orderBy, p)
	if err != nil {
		return page.Result[Product]{}, errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.productBus.Count(ctx, filter)
	if err != nil {
		return page.Result[Product]{}, errs.Newf(errs.Internal, "count: %s", err)
	}

	return page.NewResult(toAppProducts(products), total, p), nil
}

// QueryByID returns a product by its Ia.
func (a *App) QueryByID(ctx context.Context) (Product, error) {
	prd, err := ctxPck.GetProduct(ctx)
	if err != nil {
		return Product{}, errs.Newf(errs.Internal, "querybyid: %s", err)
	}

	return toAppProduct(prd), nil
}
