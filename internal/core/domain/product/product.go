package product

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/internal/core/domain/money"
	"github.com/Housiadas/backend-system/internal/core/domain/name"
	"github.com/Housiadas/backend-system/internal/core/domain/quantity"
)

const (
	ProductUpdatedEvent = "productapi-updated"
	ProductDeletedEvent = "productapi-deleted"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound     = errors.New("product not found")
	ErrUserDisabled = errors.New("user disabled")
	ErrInvalidCost  = errors.New("cost not valid")
)

// Product represents an individual product.
type Product struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        name.Name
	Cost        money.Money
	Quantity    quantity.Quantity
	DateCreated time.Time
	DateUpdated time.Time
}

// NewProduct is what we require from clients when adding a Product.
type NewProduct struct {
	UserID   uuid.UUID
	Name     name.Name
	Cost     money.Money
	Quantity quantity.Quantity
}

// UpdateProduct defines what information may be provided to modify an
// existing Product. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields, so we can differentiate
// between a field that was not provided and a field provided as
// explicitly blank. Normally we do not want to use pointers to basic usecase, but
// we make exceptions around marshaling/unmarshalling.
type UpdateProduct struct {
	Name     *name.Name
	Cost     *money.Money
	Quantity *quantity.Quantity
}

// QueryFilter holds the available fields a query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type QueryFilter struct {
	ID       *uuid.UUID
	Name     *name.Name
	Cost     *float64
	Quantity *int
}
