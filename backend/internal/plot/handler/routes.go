package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes defines the routes for the plot module.
func (h *PlotHandler) Routes() http.Handler {
	r := chi.NewRouter()

	// The main routing logic is now in api/server.go.
	// This function can be used to mount simple, non-nested routes if needed in the future.

	return r
}