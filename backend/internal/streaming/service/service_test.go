package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	cameraModels "github.com/rendley/vegshare/backend/internal/camera/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockCameraService struct {
	mock.Mock
}

func (m *MockCameraService) CreateCamera(ctx context.Context, name, rtspPathName string, plotID uuid.UUID) (*cameraModels.Camera, error) {
	args := m.Called(ctx, name, rtspPathName, plotID)
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

func TestStreamingService(t *testing.T) {
	ctx := context.Background()
	mockCameraSvc := new(MockCameraService)
	streamingSvc := NewService(mockCameraSvc)

	t.Run("AuthorizeStream - Success", func(t *testing.T) {
		cameraID := uuid.New()
		mockCameraSvc.On("GetCameraByID", ctx, cameraID).Return(&cameraModels.Camera{}, nil).Once()

		ok, err := streamingSvc.AuthorizeStream(ctx, cameraID)

		assert.True(t, ok)
		assert.NoError(t, err)
		mockCameraSvc.AssertExpectations(t)
	})

	t.Run("AuthorizeStream - Camera Not Found", func(t *testing.T) {
		cameraID := uuid.New()
		mockCameraSvc.On("GetCameraByID", ctx, cameraID).Return(nil, errors.New("not found")).Once()

		ok, err := streamingSvc.AuthorizeStream(ctx, cameraID)

		assert.False(t, ok)
		assert.Error(t, err)
		mockCameraSvc.AssertExpectations(t)
	})
}
