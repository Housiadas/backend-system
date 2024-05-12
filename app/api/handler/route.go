package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Housiadas/backend-system/business/auth"
)

// Routes returns applications router
func (h *Handler) Routes() *chi.Mux {
	mid := h.Web.Mid

	authenticate := mid.Bearer()
	ruleAny := mid.Authorize(auth.RuleAny)
	ruleUserOnly := mid.Authorize(auth.RuleUserOnly)
	ruleAdmin := mid.Authorize(auth.RuleAdminOnly)

	ruleAuthorizeUser := mid.AuthorizeUser(auth.RuleAdminOrSubject)
	ruleAuthorizeAdmin := mid.AuthorizeUser(auth.RuleAdminOnly)
	ruleAuthorizeProduct := mid.AuthorizeProduct(auth.RuleAdminOrSubject)

	apiRouter := chi.NewRouter()
	apiRouter.Use(
		mid.RequestID,
		middleware.SetHeader("Content-Type", "application/json"),
		middleware.GetHead,
		mid.Recoverer(),
	)

	apiRouter.Route("/v1", func(r chi.Router) {
		r.Post("/authapi/authenticate", h.Web.Respond.Respond(h.authenticate))
		r.With(authenticate).Get("/authapi/authorize", h.Web.Respond.Respond(h.authorize))

		// Users
		r.With(authenticate).Route("/users", func(u chi.Router) {
			u.With(ruleAuthorizeAdmin).Get("/", h.Web.Respond.Respond(h.userQuery))
			u.With(ruleAuthorizeUser).Get("/{user_id}", h.Web.Respond.Respond(h.userQueryByID))
			u.With(ruleAdmin).Post("/users", h.Web.Respond.Respond(h.userCreate))
			u.With(ruleAuthorizeAdmin).Put("/role/{user_id}", h.Web.Respond.Respond(h.updateRole))
			u.With(ruleAuthorizeUser).Put("/{user_id}", h.Web.Respond.Respond(h.userUpdate))
			u.With(ruleAuthorizeUser).Delete("/{user_id}", h.Web.Respond.Respond(h.userDelete))
		})

		// Products
		r.With(authenticate).Route("/products", func(p chi.Router) {
			p.With(ruleAny).Get("/", h.Web.Respond.Respond(h.productQuery))
			p.With(ruleUserOnly).Post("/", h.Web.Respond.Respond(h.productCreate))
			p.With(ruleAuthorizeProduct).Get("/{product_id}", h.Web.Respond.Respond(h.productQueryByID))
			p.With(ruleAuthorizeProduct).Put("/{product_id}", h.Web.Respond.Respond(h.productUpdate))
			p.With(ruleAuthorizeProduct).Delete("/{product_id}", h.Web.Respond.Respond(h.productDelete))
		})
	})

	// System Routes
	router := chi.NewRouter()
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		h.Web.Respond.Respond(h.notFound)
	})
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		h.Web.Respond.Respond(h.notAllowed)
	})
	router.Get("/readiness", h.Web.Respond.Respond(h.readiness))
	router.Get("/liveness", h.Web.Respond.Respond(h.liveness))
	router.Handle("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("./doc.json"),
	))
	router.Get("/swagger/doc.json", h.Swagger)

	router.Mount("/api", apiRouter)
	return router
}
