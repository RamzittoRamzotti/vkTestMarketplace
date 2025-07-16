package router

import (
	"github.com/go-chi/chi/v5"
	"vkTestMarketplace/internal/http-server/handlers"
	"vkTestMarketplace/internal/http-server/middleware/auth"
)

func NewRouter(authHandlers *handlers.AuthHandlers, adHandlers *handlers.AdHandlers, jwtSecret string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(auth.OptionalAuthMiddleware(jwtSecret))

	r.Post("/register", authHandlers.RegisterHandler)
	r.Post("/login", authHandlers.LoginHandler)

	adHandlerWithMiddleware := &auth.AdHandler{Handlers: adHandlers}
	r.Group(func(r chi.Router) {
		r.Use(adHandlerWithMiddleware.AuthMiddleware)
		r.Post("/ads", adHandlers.CreateAdHandler)
	})

	r.Get("/ads", adHandlers.ListAdsHandler)

	return r
}
