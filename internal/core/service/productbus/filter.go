package productbus

import (
	"github.com/Housiadas/backend-system/internal/core/domain/name"
	"github.com/google/uuid"
)

// QueryFilter holds the available fields a query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type QueryFilter struct {
	ID       *uuid.UUID
	Name     *name.Name
	Cost     *float64
	Quantity *int
}
