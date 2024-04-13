package userapi

import (
	"net/http"

	"github.com/Housiadas/simple-banking-system/app/domain/userapp"
	"github.com/Housiadas/simple-banking-system/business/domain/userbus"
	"github.com/Housiadas/simple-banking-system/foundation/logger"
	"github.com/Housiadas/simple-banking-system/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log     *logger.Logger
	UserBus *userbus.Core
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	//authen := mid.Authenticate(cfg.Log, cfg.AuthSrv)
	//ruleAdmin := mid.Authorize(cfg.Log, cfg.AuthSrv, auth.RuleAdminOnly)
	//ruleAuthorizeUser := mid.AuthorizeUser(cfg.Log, cfg.AuthSrv, cfg.UserBus, auth.RuleAdminOrSubject)
	//ruleAuthorizeAdmin := mid.AuthorizeUser(cfg.Log, cfg.AuthSrv, cfg.UserBus, auth.RuleAdminOnly)

	api := newAPI(userapp.NewCore(cfg.UserBus))
	//app.Handle(http.MethodGet, version, "/users", api.query)
	//app.Handle(http.MethodGet, version, "/users/{user_id}", api.queryByID)
	app.Handle(http.MethodPost, version, "/users", api.create)
	//app.Handle(http.MethodPut, version, "/users/role/{user_id}", api.updateRole)
	//app.Handle(http.MethodPut, version, "/users/{user_id}", api.update)
	//app.Handle(http.MethodDelete, version, "/users/{user_id}", api.delete)
}
