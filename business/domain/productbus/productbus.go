// Package productbus provides business access to product domain.
package productbus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/business/data/transaction"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/sys/order"
	"github.com/Housiadas/backend-system/foundation/kafka"
	"github.com/Housiadas/backend-system/foundation/logger"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound     = errors.New("product not found")
	ErrUserDisabled = errors.New("user disabled")
	ErrInvalidCost  = errors.New("cost not valid")
)

// Storer interface declares the behavior this package needs to persists and retrieve data.
type Storer interface {
	ExecuteUnderTransaction(tx transaction.Transaction) (Storer, error)
	Create(ctx context.Context, prd Product) error
	Update(ctx context.Context, prd Product) error
	Delete(ctx context.Context, prd Product) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Product, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, productID uuid.UUID) (Product, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Product, error)
}

// Business manages the set of APIs for product access.
type Business struct {
	log      *logger.Logger
	storer   Storer
	userBus  *userbus.Business
	producer *kafka.ProducerClient
}

// NewBusiness constructs a product core API for use.
func NewBusiness(log *logger.Logger, storer Storer, userBus *userbus.Business, producer *kafka.ProducerClient) *Business {
	return &Business{
		log:      log,
		storer:   storer,
		userBus:  userBus,
		producer: producer,
	}
}

// ExecuteUnderTransaction constructs a new Business value that will use the
// specified transaction in any store related calls.
func (c *Business) ExecuteUnderTransaction(tx transaction.Transaction) (*Business, error) {
	storer, err := c.storer.ExecuteUnderTransaction(tx)
	if err != nil {
		return nil, err
	}

	userBus, err := c.userBus.ExecuteUnderTransaction(tx)
	if err != nil {
		return nil, err
	}

	core := Business{
		log:      c.log,
		storer:   storer,
		userBus:  userBus,
		producer: c.producer,
	}

	return &core, nil
}

// Create adds a new product to the system.
func (c *Business) Create(ctx context.Context, np NewProduct) (Product, error) {
	usr, err := c.userBus.QueryByID(ctx, np.UserID)
	if err != nil {
		return Product{}, fmt.Errorf("user.querybyid: %s: %w", np.UserID, err)
	}

	if np.Cost < 0 {
		return Product{}, ErrInvalidCost
	}

	if !usr.Enabled {
		return Product{}, ErrUserDisabled
	}

	now := time.Now()

	prd := Product{
		ID:          uuid.New(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		UserID:      np.UserID,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.storer.Create(ctx, prd); err != nil {
		return Product{}, fmt.Errorf("create: %w", err)
	}

	return prd, nil
}

// Update modifies information about a product.
func (c *Business) Update(ctx context.Context, prd Product, up UpdateProduct) (Product, error) {
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
		return Product{}, fmt.Errorf("update: %w", err)
	}

	return prd, nil
}

// Delete removes the specified product.
func (c *Business) Delete(ctx context.Context, prd Product) error {
	if err := c.storer.Delete(ctx, prd); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing products.
func (c *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Product, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	prds, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}

// Count returns the total number of products.
func (c *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	if err := filter.Validate(); err != nil {
		return 0, err
	}

	return c.storer.Count(ctx, filter)
}

// QueryByID finds the product by the specified ID.
func (c *Business) QueryByID(ctx context.Context, productID uuid.UUID) (Product, error) {
	prd, err := c.storer.QueryByID(ctx, productID)
	if err != nil {
		return Product{}, fmt.Errorf("query: productID[%s]: %w", productID, err)
	}

	return prd, nil
}

// QueryByUserID finds the products by a specified User ID.
func (c *Business) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Product, error) {
	prds, err := c.storer.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}
