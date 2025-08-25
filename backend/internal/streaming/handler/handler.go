package handler

import (
	"net/http"

	"github.com/rendley/vegshare/backend/internal/streaming/service"
	"github.com/sirupsen/logrus"
)

// StreamingHandler handles WebSocket connections for video streams.
type StreamingHandler struct {
	service service.Service
	logger  *logrus.Logger
}

// NewStreamingHandler creates a new StreamingHandler.
func NewStreamingHandler(s service.Service, l *logrus.Logger) *StreamingHandler {
	return &StreamingHandler{
		service: s,
		logger:  l,
	}
}

// handleWebSocket is the actual handler function for the WebSocket endpoint.
func (h *StreamingHandler) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// The service now handles everything, from upgrading the connection to proxying.
	h.service.HandleStream(w, r)
}
