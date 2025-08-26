package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/rendley/vegshare/backend/internal/camera/service"
	"github.com/rendley/vegshare/backend/pkg/config"
	"github.com/rendley/vegshare/backend/pkg/middleware"
	"github.com/sirupsen/logrus"
)

// Service handles the proxying of WebSocket connections to the media server.
type Service interface {
	HandleStream(w http.ResponseWriter, r *http.Request)
	ProxyHLS(ctx context.Context, w http.ResponseWriter, r *http.Request, cameraPath, fileName string)
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

// ProxyHLS fetches HLS files from mediamtx and streams them to the client.
func (s *serviceImpl) ProxyHLS(ctx context.Context, w http.ResponseWriter, r *http.Request, cameraPath, fileName string) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		s.log.Error("Failed to get user ID from context for HLS proxy")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// TODO: Implement authorization logic here.
	// 1. Get camera by path `cameraPath`.
	// 2. Get plot from camera.plot_id.
	// 3. Check if `userID` has an active lease for that plot.
	// 4. If not, return http.StatusForbidden.
	s.log.Infof("User %s authorized to access HLS stream for %s", userID, cameraPath)

	// Construct the target URL to mediamtx's HLS server
	// Note: mediamtx HLS is on port 8888 by default in our config
	hlsURL := fmt.Sprintf("http://%s:8888/%s/%s", s.cfg.MediaMTX.Host, cameraPath, fileName)
	s.log.Infof("Proxying HLS request to: %s", hlsURL)

	// Create a new request to the target URL
	req, err := http.NewRequestWithContext(ctx, "GET", hlsURL, nil)
	if err != nil {
		s.log.Errorf("Failed to create HLS proxy request: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Copy headers from the original request
	req.Header.Set("User-Agent", r.Header.Get("User-Agent"))

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.log.Errorf("Failed to fetch HLS content from mediamtx: %v", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy headers from the mediamtx response to our response
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set the status code
	w.WriteHeader(resp.StatusCode)

	// Copy the body (the .m3u8 or .ts file content)
	io.Copy(w, resp.Body)
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
	// Copy essential headers from the original request
	if origin := r.Header.Get("Origin"); origin != "" {
		requestHeader.Set("Origin", origin)
	}
	if userAgent := r.Header.Get("User-Agent"); userAgent != "" {
		requestHeader.Set("User-Agent", userAgent)
	}
	// Set X-Forwarded-For header
	if clientIP := r.RemoteAddr; clientIP != "" {
		// If there are multiple proxies, they may append IPs. We take the first one.
		if fwd, ok := r.Header["X-Forwarded-For"]; ok {
			clientIP = fwd[0]
		}
		// Remove port if present
		if strings.Contains(clientIP, ":") {
			clientIP = strings.Split(clientIP, ":")[0]
		}
		requestHeader.Set("X-Forwarded-For", clientIP)
	}

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
