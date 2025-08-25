package service

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/rendley/vegshare/backend/internal/camera/service"
	"github.com/rendley/vegshare/backend/pkg/config"
	"github.com/sirupsen/logrus"
)

// Service handles the proxying of WebSocket connections to the media server.
type Service interface {
	HandleStream(w http.ResponseWriter, r *http.Request)
}

type serviceImpl struct {
	cfg       *config.Config
	log       *logrus.Logger
	cameraSvc service.Service
	upgrader  websocket.Upgrader
}

// NewService creates a new streaming service.
func NewService(cfg *config.Config, log *logrus.Logger, cameraSvc service.Service) Service {
	return &serviceImpl{
		cfg:       cfg,
		log:       log,
		cameraSvc: cameraSvc,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow all connections for now
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

// HandleStream upgrades the client connection to WebSocket and proxies it to mediamtx.
func (s *serviceImpl) HandleStream(w http.ResponseWriter, r *http.Request) {
	cameraPath := chi.URLParam(r, "cameraPath")
	if cameraPath == "" {
		s.log.Error("Camera path is empty")
		http.Error(w, "Camera path is required", http.StatusBadRequest)
		return
	}
	s.log.Infof("Attempting to stream camera with path: %s", cameraPath)

	// Prepare headers for the connection to the media server.
	requestHeader := http.Header{}
	// Let's only forward the User-Agent and see if that helps.
	requestHeader.Add("User-Agent", r.Header.Get("User-Agent"))

	clientConn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Errorf("Failed to upgrade client connection: %v", err)
		return
	}
	defer clientConn.Close()
	s.log.Info("Client connection upgraded to WebSocket")

	mediaServerHost := fmt.Sprintf("%s:%s", s.cfg.MediaMTX.Host, s.cfg.MediaMTX.Port)
	targetURL := url.URL{Scheme: "ws", Host: mediaServerHost, Path: "/" + cameraPath}
	s.log.Infof("Connecting to media server at: %s with headers: %v", targetURL.String(), requestHeader)

	// Connect to the media server with the prepared headers
	mediaServerConn, _, err := websocket.DefaultDialer.Dial(targetURL.String(), requestHeader)
	if err != nil {
		s.log.Errorf("Failed to connect to media server: %v", err)
		clientConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "could not connect to media server"))
		return
	}
	defer mediaServerConn.Close()
	s.log.Info("Successfully connected to media server")

	var wg sync.WaitGroup
	wg.Add(2)

	go s.proxyMessages(clientConn, mediaServerConn, &wg, "Client to Server")
	go s.proxyMessages(mediaServerConn, clientConn, &wg, "Server to Client")

	wg.Wait()
	s.log.Info("Streaming session finished")
}

// proxyMessages reads messages from the source and writes them to the destination.
func (s *serviceImpl) proxyMessages(source, dest *websocket.Conn, wg *sync.WaitGroup, direction string) {
	defer wg.Done()
	for {
		msgType, msg, err := source.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.log.Errorf("Error reading message (%s): %v", direction, err)
			} else {
				s.log.Infof("Connection closed (%s): %v", direction, err)
			}
			dest.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			break
		}

		if err := dest.WriteMessage(msgType, msg); err != nil {
			s.log.Errorf("Error writing message (%s): %v", direction, err)
			break
		}
	}
}
