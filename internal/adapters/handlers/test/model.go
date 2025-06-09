package test

import (
	"github.com/Housiadas/backend-system/internal/core/service/productbus"
	"github.com/Housiadas/backend-system/internal/core/service/userbus"
)

// User extends the dbtest user for api test support.
type User struct {
	userbus.User
	Products []productbus.Product
	Token    string
}

// SeedData represents users for api tests.
type SeedData struct {
	Users  []User
	Admins []User
}

// Table represent fields needed for running an api test.
type Table struct {
	Name       string
	URL        string
	Token      string
	Method     string
	StatusCode int
	Input      any
	GotResp    any
	ExpResp    any
	CmpFunc    func(got any, exp any) string
}
