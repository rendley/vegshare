package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes returns an empty router for the camera handler.
// The routes are currently defined in the central server.go file to handle nested routing.
func (h *CameraHandler) Routes() http.Handler {
	r := chi.NewRouter()
	return r
}
