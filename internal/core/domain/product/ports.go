package product

import (
	"context"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/pkg/order"
	"github.com/Housiadas/backend-system/pkg/page"
	"github.com/Housiadas/backend-system/pkg/pgsql"
)

// Storer interface declares the behavior this package needs to priests and retrieve data.
type Storer interface {
	NewWithTx(tx pgsql.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, prd Product) error
	Update(ctx context.Context, prd Product) error
	Delete(ctx context.Context, prd Product) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Product, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, productID uuid.UUID) (Product, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Product, error)
}
