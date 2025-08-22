package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes returns a new router for the streaming handler.
func (h *StreamingHandler) Routes() http.Handler {
	r := chi.NewRouter()

	// This will be mounted under /stream, so the final path will be /stream/{cameraID}
	r.Get("/{cameraID}", h.HandleStream)

	return r
}
