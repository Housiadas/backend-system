// Package all binds all the routes into the specified app.
package all

import (
	"github.com/Housiadas/simple-banking-system/app/api/mux"
	"github.com/Housiadas/simple-banking-system/app/api/route/sys"
	"github.com/Housiadas/simple-banking-system/app/api/route/userapi"
	"github.com/Housiadas/simple-banking-system/foundation/web"
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
	sys.Routes(app, sys.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	userapi.Routes(app, userapi.Config{
		Log:     cfg.Log,
		UserBus: cfg.BusDomain.User,
	})
}
