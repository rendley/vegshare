package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes returns a new router for the camera handler.
func (h *CameraHandler) Routes() http.Handler {
	r := chi.NewRouter()

	// The main routing logic is now in api/server.go.

	return r
}