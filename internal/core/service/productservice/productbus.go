// Package productservice provides internal access to the product core.
package productservice

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/internal/core/domain/product"
	"github.com/Housiadas/backend-system/internal/core/service/userservice"
	"github.com/Housiadas/backend-system/pkg/logger"
	"github.com/Housiadas/backend-system/pkg/order"
	"github.com/Housiadas/backend-system/pkg/otel"
	"github.com/Housiadas/backend-system/pkg/page"
	"github.com/Housiadas/backend-system/pkg/sqldb"
)

// Business manages the set of APIs for product access.
type Business struct {
	log     *logger.Logger
	userBus *userservice.Service
	storer  product.Storer
}

// NewBusiness constructs a product internal API for use.
func NewBusiness(
	log *logger.Logger,
	userBus *userservice.Service,
	storer product.Storer,
) *Business {
	b := Business{
		log:     log,
		userBus: userBus,
		storer:  storer,
	}

	return &b
}

// NewWithTx constructs a new internal value that will use the
// specified transaction in any store-related calls.
func (b *Business) NewWithTx(tx sqldb.CommitRollbacker) (*Business, error) {
	storer, err := b.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	userBus, err := b.userBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	bus := Business{
		log:     b.log,
		userBus: userBus,
		storer:  storer,
	}

	return &bus, nil
}

// Create adds a new product to the system.
func (b *Business) Create(ctx context.Context, np product.NewProduct) (product.Product, error) {
	ctx, span := otel.AddSpan(ctx, "internal.productservice.create")
	defer span.End()

	usr, err := b.userBus.QueryByID(ctx, np.UserID)
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

	if err := b.storer.Create(ctx, prd); err != nil {
		return product.Product{}, fmt.Errorf("create: %w", err)
	}

	return prd, nil
}

// Update modifies information about a product.
func (b *Business) Update(ctx context.Context, prd product.Product, up product.UpdateProduct) (product.Product, error) {
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

	if err := b.storer.Update(ctx, prd); err != nil {
		return product.Product{}, fmt.Errorf("update: %w", err)
	}

	return prd, nil
}

// Delete removes the specified product.
func (b *Business) Delete(ctx context.Context, prd product.Product) error {
	if err := b.storer.Delete(ctx, prd); err != nil {
		return fmt.Errorf("deleteUser: %w", err)
	}

	return nil
}

// Query retrieves a list of existing products.
func (b *Business) Query(ctx context.Context, filter product.QueryFilter, orderBy order.By, page page.Page) ([]product.Product, error) {
	prds, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}

// Count returns the total number of products.
func (b *Business) Count(ctx context.Context, filter product.QueryFilter) (int, error) {
	return b.storer.Count(ctx, filter)
}

// QueryByID finds the product by the specified Ib.
func (b *Business) QueryByID(ctx context.Context, productID uuid.UUID) (product.Product, error) {
	prd, err := b.storer.QueryByID(ctx, productID)
	if err != nil {
		return product.Product{}, fmt.Errorf("query: productID[%s]: %w", productID, err)
	}

	return prd, nil
}

// QueryByUserID finds the products by a specified User Ib.
func (b *Business) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]product.Product, error) {
	prds, err := b.storer.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}
