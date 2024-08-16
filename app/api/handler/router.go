package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/riandyrn/otelchi"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Housiadas/backend-system/business/domain/authbus"
)

// Routes returns applications router
func (h *Handler) Routes() *chi.Mux {
	mid := h.Web.Mid

	authenticate := mid.Bearer()
	ruleAuthorizeAny := mid.AuthorizeUser(authbus.RuleAny)
	ruleAuthorizeUserOnly := mid.AuthorizeUser(authbus.RuleUserOnly)
	ruleAuthorizeUser := mid.AuthorizeUser(authbus.RuleAdminOrSubject)
	ruleAuthorizeAdmin := mid.AuthorizeUser(authbus.RuleAdminOnly)
	ruleAuthorizeProduct := mid.AuthorizeProduct(authbus.RuleAdminOrSubject)

	tran := mid.BeginCommitRollback()

	apiRouter := chi.NewRouter()
	apiRouter.Use(
		mid.Recoverer(),
		mid.RequestID,
		mid.Otel(),
		mid.Logger(),
		middleware.SetHeader("Content-Type", "application/json"),
		middleware.GetHead,
	)

	// v1 routes
	apiRouter.Route("/v1", func(v1 chi.Router) {
		v1.Use(otelchi.Middleware(h.ServiceName, otelchi.WithChiRoutes(v1)))
		v1.Use(mid.ApiVersion("v1"))
		v1.Use(cors.Handler(cors.Options{
			AllowedOrigins: h.Cors.AllowedOrigins,
			AllowedMethods: h.Cors.AllowedMethods,
			AllowedHeaders: h.Cors.AllowedHeaders,
			ExposedHeaders: h.Cors.ExposedHeaders,
			MaxAge:         h.Cors.MaxAge,
		}))

		// Auth
		v1.Post("/auth/authenticate", h.Web.Res.Respond(h.authenticate))
		v1.With(authenticate).Get("/auth/authorize", h.Web.Res.Respond(h.authorize))

		// Users
		v1.With(authenticate).Route("/users", func(u chi.Router) {
			u.With(ruleAuthorizeAdmin).Get("/", h.Web.Res.Respond(h.userQuery))
			u.With(ruleAuthorizeAdmin).Post("/", h.Web.Res.Respond(h.userCreate))
			u.With(ruleAuthorizeUser).Get("/{user_id}", h.Web.Res.Respond(h.userQueryByID))
			u.With(ruleAuthorizeAdmin).Put("/role/{user_id}", h.Web.Res.Respond(h.updateRole))
			u.With(ruleAuthorizeUser).Put("/{user_id}", h.Web.Res.Respond(h.userUpdate))
			u.With(ruleAuthorizeUser).Delete("/{user_id}", h.Web.Res.Respond(h.userDelete))
		})

		// Products
		v1.With(authenticate).Route("/products", func(p chi.Router) {
			p.With(ruleAuthorizeAny).Get("/", h.Web.Res.Respond(h.productQuery))
			p.With(ruleAuthorizeUserOnly).Post("/", h.Web.Res.Respond(h.productCreate))
			p.With(ruleAuthorizeProduct).Get("/{product_id}", h.Web.Res.Respond(h.productQueryByID))
			p.With(ruleAuthorizeProduct).Put("/{product_id}", h.Web.Res.Respond(h.productUpdate))
			p.With(ruleAuthorizeProduct).Delete("/{product_id}", h.Web.Res.Respond(h.productDelete))
		})

		// Transaction example
		v1.With(tran).Post("/transaction", h.Web.Res.Respond(h.transaction))
	})

	// System Routes
	router := chi.NewRouter()
	router.Get("/readiness", h.Web.Res.Respond(h.readiness))
	router.Get("/liveness", h.Web.Res.Respond(h.liveness))
	router.Handle("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("./doc.json"),
	))
	router.Get("/swagger/doc.json", h.Swagger)

	router.Mount("/api", apiRouter)
	return router
}
