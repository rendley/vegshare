package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes returns a new router for the farm handler.
func (h *FarmHandler) Routes() http.Handler {
	r := chi.NewRouter()

	// All farm-related routes are now defined in the main api/server.go
	// to avoid circular dependencies and centralize routing logic.

	// Example of what could be here if routes were simple:
	// r.Get("/regions", h.GetAllRegions)

	return r
}