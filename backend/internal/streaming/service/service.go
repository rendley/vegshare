package service

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	cameraService "github.com/rendley/vegshare/backend/internal/camera/service"
	"github.com/rendley/vegshare/backend/pkg/config"
	"github.com/sirupsen/logrus"
)

// Service defines the contract for the streaming service.
type Service interface {
	HandleStream(ctx context.Context, clientConn *websocket.Conn, cameraIDStr string) error
}

// service implements the Service interface.
type service struct {
	cfg       *config.Config
	logger    *logrus.Logger
	cameraSvc cameraService.Service
}

// NewService is a constructor for the streaming service.
func NewService(cfg *config.Config, logger *logrus.Logger, cameraSvc cameraService.Service) Service {
	return &service{cfg: cfg, logger: logger, cameraSvc: cameraSvc}
}

func (s *service) HandleStream(ctx context.Context, clientConn *websocket.Conn, cameraIDStr string) error {
	cameraID, err := uuid.Parse(cameraIDStr)
	if err != nil {
		return fmt.Errorf("invalid camera ID: %w", err)
	}

	cam, err := s.cameraSvc.GetCameraByID(ctx, cameraID)
	if err != nil {
		return fmt.Errorf("failed to get camera details: %w", err)
	}

	// Construct the URL to mediamtx
	// Example: http://mediamtx:8888/ws/cam.rtsp_url
	u := url.URL{
		Scheme: "ws",
		Host:   s.cfg.MediaMTX.Host + ":" + s.cfg.MediaMTX.Port,
		Path:   "/" + cameraIDStr, // Use camera ID as the path
	}

	// We need to re-encode the original RTSP url to pass it as a query param to mediamtx
	q := u.Query()
	q.Set("url", cam.RTSPURL)
	u.RawQuery = q.Encode()

	s.logger.Infof("Connecting to mediamtx: %s", u.String())

	// Connect to mediamtx
	mtxConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to mediamtx: %w", err)
	}
	defer mtxConn.Close()

	s.logger.Info("Successfully connected to mediamtx")

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine to forward messages from client to mediamtx
	go func() {
		defer wg.Done()
		for {
			messageType, p, err := clientConn.ReadMessage()
			if err != nil {
				s.logger.Debugf("Client read error: %v", err)
				mtxConn.Close() // Close connection to mediamtx if client disconnects
				return
			}
			if err := mtxConn.WriteMessage(messageType, p); err != nil {
				s.logger.Debugf("Mediamtx write error: %v", err)
				return
			}
		}
	}()

	// Goroutine to forward messages from mediamtx to client
	go func() {
		defer wg.Done()
		for {
			messageType, p, err := mtxConn.ReadMessage()
			if err != nil {
				s.logger.Debugf("Mediamtx read error: %v", err)
				clientConn.Close() // Close connection to client if mediamtx disconnects
				return
			}
			if err := clientConn.WriteMessage(messageType, p); err != nil {
				s.logger.Debugf("Client write error: %v", err)
				return
			}
		}
	}()

	wg.Wait()
	return nil
}