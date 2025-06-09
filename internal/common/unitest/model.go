package unitest

import (
	"context"

	"github.com/Housiadas/backend-system/internal/core/service/productbus"
	"github.com/Housiadas/backend-system/internal/core/service/userbus"
)

// User represents an cli user specified for the test.
type User struct {
	userbus.User
	Products []productbus.Product
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
