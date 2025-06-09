package test

import (
	"github.com/Housiadas/backend-system/internal/core/service/productservice"
	"github.com/Housiadas/backend-system/internal/core/service/userservice"
)

// User extends the dbtest user for api test support.
type User struct {
	userservice.User
	Products []productservice.Product
}

// SeedData represents users for api tests.
type SeedData struct {
	Users  []User
	Admins []User
}

// Table represent fields needed for running an api test.
type Table struct {
	Name    string
	Path    string
	Input   any
	GotResp any
	ExpResp any
	CmpFunc func(got any, exp any) string
}
