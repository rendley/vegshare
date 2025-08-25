package service

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	cameraModels "github.com/rendley/vegshare/backend/internal/camera/models"
	"github.com/rendley/vegshare/backend/pkg/config"
	"github.com/rendley/vegshare/backend/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockCameraService struct {
	mock.Mock
}

func (m *MockCameraService) CreateCamera(ctx context.Context, name, rtspURL string, plotID uuid.UUID) (*cameraModels.Camera, error) {
	args := m.Called(ctx, name, rtspURL, plotID)
	return args.Get(0).(*cameraModels.Camera), args.Error(1)
}

func (m *MockCameraService) GetCamerasByPlotID(ctx context.Context, plotID uuid.UUID) ([]cameraModels.Camera, error) {
	args := m.Called(ctx, plotID)
	return args.Get(0).([]cameraModels.Camera), args.Error(1)
}

func (m *MockCameraService) GetCameraByID(ctx context.Context, cameraID uuid.UUID) (*cameraModels.Camera, error) {
	args := m.Called(ctx, cameraID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cameraModels.Camera), args.Error(1)
}

func (m *MockCameraService) DeleteCamera(ctx context.Context, cameraID uuid.UUID) error {
	args := m.Called(ctx, cameraID)
	return args.Error(0)
}

// --- Tests ---

func TestStreamingService_HandleStream(t *testing.T) {
	// Setup
	log := logger.New() // Use real logger for simplicity in tests
	cfg := &config.Config{} // Empty config, as we are not connecting to a real mediamtx
	mockCameraSvc := new(MockCameraService)
	streamingSvc := NewService(cfg, log, mockCameraSvc)

	// Create a test server
	server := httptest.NewServer(nil) // We don't need a real handler for this test
	defer server.Close()

	// This is a dummy dialer, we won't actually connect
	dialer := websocket.Dialer{}

	t.Run("should return error if camera not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		cameraID := uuid.New()
		wsURL := "ws://" + server.Listener.Addr().String() + "/stream/" + cameraID.String()

		// We don't actually need a real websocket connection for this test case,
		// as the error happens before the connection is even used.
		conn, _, _ := dialer.Dial(wsURL, nil)

		mockCameraSvc.On("GetCameraByID", ctx, cameraID).Return(nil, errors.New("not found")).Once()

		// Act
		err := streamingSvc.HandleStream(ctx, conn, cameraID.String())

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get camera details")
		mockCameraSvc.AssertExpectations(t)
	})
}