package handler

import "github.com/go-chi/chi/v5"

// RegisterRoutes registers the auth routes.
func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/register", h.Register)
	r.Post("/api/v1/login", h.Login)
}
