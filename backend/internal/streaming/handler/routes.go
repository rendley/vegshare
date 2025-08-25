package handler

import "github.com/go-chi/chi/v5"

// Routes returns the routing tree for the streaming service.
func (h *StreamingHandler) Routes() chi.Router {
	r := chi.NewRouter()
	// The path now includes a wildcard to capture the camera path.
	r.Get("/ws/{cameraPath:*}", h.handleWebSocket)
	return r
}