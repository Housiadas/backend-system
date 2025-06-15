// Package productcore provides internal access to the product core.
package productcore

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/internal/core/domain/product"
	"github.com/Housiadas/backend-system/internal/core/service/usercore"
	"github.com/Housiadas/backend-system/pkg/logger"
	"github.com/Housiadas/backend-system/pkg/order"
	"github.com/Housiadas/backend-system/pkg/otel"
	"github.com/Housiadas/backend-system/pkg/page"
	"github.com/Housiadas/backend-system/pkg/sqldb"
)

// Core manages the set of APIs for product access.
type Core struct {
	log     *logger.Logger
	userBus *usercore.Core
	storer  product.Storer
}

// NewCore constructs a product internal API for use.
func NewCore(
	log *logger.Logger,
	userBus *usercore.Core,
	storer product.Storer,
) *Core {
	b := Core{
		log:     log,
		userBus: userBus,
		storer:  storer,
	}

	return &b
}

// NewWithTx constructs a new internal value that will use the
// specified transaction in any store-related calls.
func (c *Core) NewWithTx(tx sqldb.CommitRollbacker) (*Core, error) {
	storer, err := c.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	userBus, err := c.userBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	bus := Core{
		log:     c.log,
		userBus: userBus,
		storer:  storer,
	}

	return &bus, nil
}

// Create adds a new product to the system.
func (c *Core) Create(ctx context.Context, np product.NewProduct) (product.Product, error) {
	ctx, span := otel.AddSpan(ctx, "internal.productcore.create")
	defer span.End()

	usr, err := c.userBus.QueryByID(ctx, np.UserID)
	if err != nil {
		return product.Product{}, fmt.Errorf("user.querybyid: %s: %w", np.UserID, err)
	}

	if !usr.Enabled {
		return product.Product{}, product.ErrUserDisabled
	}

	now := time.Now()

	prd := product.Product{
		ID:          uuid.New(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		UserID:      np.UserID,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.storer.Create(ctx, prd); err != nil {
		return product.Product{}, fmt.Errorf("create: %w", err)
	}

	return prd, nil
}

// Update modifies information about a product.
func (c *Core) Update(ctx context.Context, prd product.Product, up product.UpdateProduct) (product.Product, error) {
	if up.Name != nil {
		prd.Name = *up.Name
	}

	if up.Cost != nil {
		prd.Cost = *up.Cost
	}

	if up.Quantity != nil {
		prd.Quantity = *up.Quantity
	}

	prd.DateUpdated = time.Now()

	if err := c.storer.Update(ctx, prd); err != nil {
		return product.Product{}, fmt.Errorf("update: %w", err)
	}

	return prd, nil
}

// Delete removes the specified product.
func (c *Core) Delete(ctx context.Context, prd product.Product) error {
	if err := c.storer.Delete(ctx, prd); err != nil {
		return fmt.Errorf("deleteUser: %w", err)
	}

	return nil
}

// Query retrieves a list of existing products.
func (c *Core) Query(ctx context.Context, filter product.QueryFilter, orderBy order.By, page page.Page) ([]product.Product, error) {
	prds, err := c.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}

// Count returns the total number of products.
func (c *Core) Count(ctx context.Context, filter product.QueryFilter) (int, error) {
	return c.storer.Count(ctx, filter)
}

// QueryByID finds the product by the specified Ib.
func (c *Core) QueryByID(ctx context.Context, productID uuid.UUID) (product.Product, error) {
	prd, err := c.storer.QueryByID(ctx, productID)
	if err != nil {
		return product.Product{}, fmt.Errorf("query: productID[%s]: %w", productID, err)
	}

	return prd, nil
}

// QueryByUserID finds the products by a specified User Ib.
func (c *Core) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]product.Product, error) {
	prds, err := c.storer.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}
