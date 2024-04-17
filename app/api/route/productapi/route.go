package productapi

import (
	"net/http"

	"github.com/Housiadas/backend-system/app/api/mid"
	"github.com/Housiadas/backend-system/app/domain/productapp"
	"github.com/Housiadas/backend-system/business/auth"
	"github.com/Housiadas/backend-system/business/domain/productbus"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/foundation/logger"
	"github.com/Housiadas/backend-system/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log        *logger.Logger
	Auth       *auth.Auth
	UserBus    *userbus.Core
	ProductBus *productbus.Core
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.Auth)
	ruleAny := mid.Authorize(cfg.Auth, auth.RuleAny)
	ruleUserOnly := mid.Authorize(cfg.Auth, auth.RuleUserOnly)
	ruleAuthorizeProduct := mid.AuthorizeProduct(cfg.Auth, cfg.ProductBus)

	api := newAPI(productapp.NewCore(cfg.ProductBus))
	app.Handle(http.MethodGet, version, "/products", api.query, authen, ruleAny)
	app.Handle(http.MethodGet, version, "/products/{product_id}", api.queryByID, authen, ruleAuthorizeProduct)
	app.Handle(http.MethodPost, version, "/products", api.create, authen, ruleUserOnly)
	app.Handle(http.MethodPut, version, "/products/{product_id}", api.update, authen, ruleAuthorizeProduct)
	app.Handle(http.MethodDelete, version, "/products/{product_id}", api.delete, authen, ruleAuthorizeProduct)
}
