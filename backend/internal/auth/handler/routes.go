
package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes returns a new router for the auth handler.
func (h *AuthHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/register", h.Register)
	r.Post("/login", h.Login)

	return r
}

