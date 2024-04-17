package route

import (
	"github.com/Housiadas/backend-system/app/api/mux"
	"github.com/Housiadas/backend-system/app/api/route/authapi"
	"github.com/Housiadas/backend-system/app/api/route/productapi"
	"github.com/Housiadas/backend-system/app/api/route/sys"
	"github.com/Housiadas/backend-system/app/api/route/userapi"
	"github.com/Housiadas/backend-system/foundation/web"
)

// Routes constructs the add value which provides the implementation of
//
//	RouteAdder for specifying what routes to bind to this instance.
func Routes() add {
	return add{}
}

type add struct{}

// Add implements the RouterAdder interface.
func (add) Add(app *web.App, cfg mux.Config) {
	// Authentication Routes
	authapi.Routes(app, authapi.Config{
		UserBus: cfg.BusDomain.User,
		Auth:    cfg.Auth,
	})

	// User Domain Routes
	userapi.Routes(app, userapi.Config{
		Log:     cfg.Log,
		Auth:    cfg.Auth,
		UserBus: cfg.BusDomain.User,
	})

	// Product Domain Routes
	productapi.Routes(app, productapi.Config{
		Log:        cfg.Log,
		Auth:       cfg.Auth,
		UserBus:    cfg.BusDomain.User,
		ProductBus: cfg.BusDomain.Product,
	})

	// System Routes
	sys.Routes(app, sys.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})
}
