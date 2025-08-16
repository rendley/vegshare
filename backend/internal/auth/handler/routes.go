package handler

import "github.com/go-chi/chi/v5"

func (h *AuthHandler) RegisterRouter(r chi.Router) {
	r.Post("/api/v1/register", h.registerHandler)
	r.Post("/api/v1/login", h.loginHandler)
}
