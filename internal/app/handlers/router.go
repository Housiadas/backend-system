package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/riandyrn/otelchi"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Housiadas/backend-system/internal/core/service/authcore"
)

// Routes returns applications router
func (h *Handler) Routes() *chi.Mux {
	mid := h.Web.Middleware

	// Bearer middleware
	authenticate := mid.Bearer()

	// authorization middleware
	ruleAny := mid.Authorize(authcore.RuleAny)
	ruleAdmin := mid.Authorize(authcore.RuleAdminOnly)
	ruleUserOnly := mid.Authorize(authcore.RuleUserOnly)

	// authorization for resource (entity) actions
	// Check if a user is allowed to modify another user's resources
	requestUserAuthorizeAdmin := mid.UserPermissions(authcore.RuleAdminOnly)
	requestUserAdminOrSubject := mid.UserPermissions(authcore.RuleAdminOrSubject)
	requestProductAdminOrSubject := mid.ProductPermissions(authcore.RuleAdminOrSubject)

	tran := mid.BeginCommitRollback()

	apiRouter := chi.NewRouter()
	apiRouter.Use(
		mid.Recoverer(),
		mid.RequestID,
		mid.Logger(),
		mid.Otel(),
		middleware.SetHeader("Content-Type", "application/json"),
		middleware.GetHead,
	)

	// v1 routes
	apiRouter.Route("/v1", func(v1 chi.Router) {
		v1.Use(
			mid.ApiVersion("v1"),
			otelchi.Middleware(h.ServiceName, otelchi.WithChiRoutes(v1)),
		)
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
			u.With(ruleAdmin).Get("/", h.Web.Res.Respond(h.userQuery))
			u.With(ruleAdmin).Post("/", h.Web.Res.Respond(h.userCreate))
			u.With(requestUserAdminOrSubject).Get("/{user_id}", h.Web.Res.Respond(h.userQueryByID))
			u.With(requestUserAuthorizeAdmin).Put("/role/{user_id}", h.Web.Res.Respond(h.updateRole))
			u.With(requestUserAdminOrSubject).Put("/{user_id}", h.Web.Res.Respond(h.userUpdate))
			u.With(requestUserAdminOrSubject).Delete("/{user_id}", h.Web.Res.Respond(h.userDelete))
		})

		// Products
		v1.With(authenticate).Route("/products", func(p chi.Router) {
			p.With(ruleAny).Get("/", h.Web.Res.Respond(h.productQuery))
			p.With(ruleUserOnly).Post("/", h.Web.Res.Respond(h.productCreate))
			p.With(requestProductAdminOrSubject).Get("/{product_id}", h.Web.Res.Respond(h.productQueryByID))
			p.With(requestProductAdminOrSubject).Put("/{product_id}", h.Web.Res.Respond(h.productUpdate))
			p.With(requestProductAdminOrSubject).Delete("/{product_id}", h.Web.Res.Respond(h.productDelete))
		})

		// Transaction example
		v1.With(tran).Post("/transaction", h.Web.Res.Respond(h.transaction))
	})

	// System Routes
	router := chi.NewRouter()
	router.Get("/readiness", h.Web.Res.Respond(h.readiness))
	router.Get("/liveness", h.Web.Res.Respond(h.liveness))
	router.Get("/swagger/doc.json", h.Swagger)
	router.Handle("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("./doc.json"),
	))

	router.Mount("/api", apiRouter)
	return router
}
