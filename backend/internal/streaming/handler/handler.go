package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rendley/vegshare/backend/internal/streaming/service"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default. In production, you should have a whitelist.
		return true
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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Errorf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	h.logger.Info("WebSocket connection established")

	// For now, just echo messages back.
	// This is where the WebRTC signaling logic will go.
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			h.logger.Errorf("Error reading message: %v", err)
			break
		}
		h.logger.Infof("Received message: %s", p)

		if err := conn.WriteMessage(messageType, p); err != nil {
			h.logger.Errorf("Error writing message: %v", err)
			break
		}
	}
}
