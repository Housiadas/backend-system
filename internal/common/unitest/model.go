package unitest

import (
	"context"

	"github.com/Housiadas/backend-system/internal/core/domain/audit"
	"github.com/Housiadas/backend-system/internal/core/domain/product"
	"github.com/Housiadas/backend-system/internal/core/domain/user"
)

// User represents a user specified for the test.
type User struct {
	user.User
	Products []product.Product
	Audits   []audit.Audit
}

// SeedData represents data seeded for the test.
type SeedData struct {
	Users  []User
	Admins []User
}

// The Table represents fields needed for running a unit test.
type Table struct {
	Name    string
	ExpResp any
	ExcFunc func(ctx context.Context) any
	CmpFunc func(got any, exp any) string
}
