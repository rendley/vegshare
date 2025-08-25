package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/rendley/vegshare/backend/internal/streaming/service"
	"github.com/rendley/vegshare/backend/pkg/middleware"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default. In production, you should have a whitelist.
		return true;
	},
}

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

// HandleStream upgrades the HTTP connection to a WebSocket and handles the stream.
func (h *StreamingHandler) HandleStream(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.logger.Error("Failed to get user ID from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get camera ID from URL
	cameraIDStr := chi.URLParam(r, "cameraID")

	// Upgrade the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Errorf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	h.logger.Infof("User %s establishing stream for camera %s", userID, cameraIDStr)

	// Pass the connection to the service layer to handle the WebRTC signaling
	if err := h.service.HandleStream(r.Context(), conn, cameraIDStr);
	 err != nil {
		h.logger.Errorf("Error during WebRTC stream handling: %v", err)
		// The connection is likely closed by the service layer, so we just log the error.
	}

	h.logger.Infof("Stream for camera %s ended", cameraIDStr)
}
