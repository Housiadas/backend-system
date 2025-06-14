package test

import (
	"github.com/Housiadas/backend-system/internal/core/domain/product"
	"github.com/Housiadas/backend-system/internal/core/domain/user"
)

// User extends the dbtest user for api test support.
type User struct {
	user.User
	Products []product.Product
}

// SeedData represents users for api tests.
type SeedData struct {
	Users  []User
	Admins []User
}

// Table represents fields needed for running an api test.
type Table struct {
	Name    string
	Path    string
	Input   any
	GotResp any
	ExpResp any
	CmpFunc func(got any, exp any) string
}
