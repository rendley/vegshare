package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes returns a new router for the camera handler.
func (h *CameraHandler) Routes() http.Handler {
	r := chi.NewRouter()

	// Note: These routes will be mounted under a specific path in the main API router.
	// For example, the GET and POST routes will be under /farm/plots/{plotID}/cameras
	// and the DELETE route will be under /cameras/{cameraID}.

	// This file just defines the sub-router for camera-specific endpoints.
	r.Get("/", h.GetCamerasByPlotID) // Mounted under /farm/plots/{plotID}/cameras
	r.Post("/", h.CreateCamera)      // Mounted under /farm/plots/{plotID}/cameras

	// This route would be mounted separately at a higher level
	// r.Delete("/{cameraID}", h.DeleteCamera) -> This needs to be handled in the main router assembly

	return r
}
