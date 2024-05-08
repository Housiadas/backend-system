package systemapi

import (
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/Housiadas/backend-system/app/domain/systemapp"
	"github.com/Housiadas/backend-system/foundation/logger"
	"github.com/Housiadas/backend-system/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {

	api := newAPI(systemapp.NewApp(cfg.Build, cfg.Log, cfg.DB))
	app.HandleNoMiddleware(http.MethodGet, "", "/readiness", api.readiness)
	app.HandleNoMiddleware(http.MethodGet, "", "/liveness", api.liveness)
	//app.HandleNoMiddleware(http.MethodGet, "", "/swagger/", api.swagger)

	mux := http.NewServeMux()
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("./swagger.json"),
	))
	app.HandleNoMiddleware(http.MethodGet, "", "/swagger/swagger.json", api.swagger)

}
