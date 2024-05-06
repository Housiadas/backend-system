package authapi

import (
	"net/http"

	"github.com/Housiadas/backend-system/app/api/mid"
	"github.com/Housiadas/backend-system/app/domain/userapp"
	"github.com/Housiadas/backend-system/business/auth"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	UserBus *userbus.Business
	Auth    *auth.Auth
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authenticate := mid.Authenticate(cfg.Auth)
	api := newAPI(userapp.NewCoreWithAuth(cfg.UserBus, cfg.Auth), cfg.Auth)

	app.Handle(http.MethodPost, version, "/auth/authenticate", api.authenticate)
	app.Handle(http.MethodGet, version, "/auth/authorize", api.authorize, authenticate)
	//app.Handle(http.MethodGet, version, "/auth/token/{kid}", api.token, authen)
}
