
package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes returns a new router for the user handler.
func (h *UserHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/me", h.GetProfile)
	r.Patch("/me", h.UpdateProfile)
	r.Delete("/me", h.DeleteAccount)

	return r
}

