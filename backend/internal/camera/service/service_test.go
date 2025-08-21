package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/camera/models"
	plotModels "github.com/rendley/vegshare/backend/internal/plot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockCameraRepository struct {
	mock.Mock
}

func (m *MockCameraRepository) CreateCamera(ctx context.Context, camera *models.Camera) error {
	args := m.Called(ctx, camera)
	return args.Error(0)
}

func (m *MockCameraRepository) GetCamerasByPlotID(ctx context.Context, plotID uuid.UUID) ([]models.Camera, error) {
	args := m.Called(ctx, plotID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Camera), args.Error(1)
}

func (m *MockCameraRepository) DeleteCamera(ctx context.Context, cameraID uuid.UUID) error {
	args := m.Called(ctx, cameraID)
	return args.Error(0)
}

// MockPlotService is a mock for the plot service
type MockPlotService struct {
	mock.Mock
}

func (m *MockPlotService) GetPlotByID(ctx context.Context, id uuid.UUID) (*plotModels.Plot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*plotModels.Plot), args.Error(1)
}

// Dummy implementations for other plot service methods to satisfy the interface
func (m *MockPlotService) CreatePlot(ctx context.Context, name, size string, greenhouseID uuid.UUID) (*plotModels.Plot, error) { return nil, nil }
func (m *MockPlotService) GetPlotsByGreenhouse(ctx context.Context, greenhouseID uuid.UUID) ([]plotModels.Plot, error) { return nil, nil }
func (m *MockPlotService) UpdatePlot(ctx context.Context, id uuid.UUID, name, size, status string) (*plotModels.Plot, error) { return nil, nil }
func (m *MockPlotService) DeletePlot(ctx context.Context, id uuid.UUID) error { return nil }


// --- Tests ---

func TestCameraService(t *testing.T) {
	ctx := context.Background()
	mockCameraRepo := new(MockCameraRepository)
	mockPlotSvc := new(MockPlotService)
	cameraSvc := NewService(mockCameraRepo, mockPlotSvc)

	t.Run("CreateCamera - Success", func(t *testing.T) {
		plotID := uuid.New()
		mockPlotSvc.On("GetPlotByID", ctx, plotID).Return(&plotModels.Plot{}, nil).Once()
		mockCameraRepo.On("CreateCamera", ctx, mock.AnythingOfType("*models.Camera")).Return(nil).Once()

		camera, err := cameraSvc.CreateCamera(ctx, "Test Cam", "test_cam", plotID)

		assert.NoError(t, err)
		assert.NotNil(t, camera)
		assert.Equal(t, "Test Cam", camera.Name)
		mockPlotSvc.AssertExpectations(t)
		mockCameraRepo.AssertExpectations(t)
	})

	t.Run("CreateCamera - Plot Not Found", func(t *testing.T) {
		plotID := uuid.New()
		mockPlotSvc.On("GetPlotByID", ctx, plotID).Return(nil, errors.New("not found")).Once()

		camera, err := cameraSvc.CreateCamera(ctx, "Test Cam", "test_cam", plotID)

		assert.Error(t, err)
		assert.Nil(t, camera)
		mockPlotSvc.AssertExpectations(t)
	})

	t.Run("GetCamerasByPlotID - Success", func(t *testing.T) {
		plotID := uuid.New()
		expectedCameras := []models.Camera{{ID: uuid.New()}, {ID: uuid.New()}}
		mockPlotSvc.On("GetPlotByID", ctx, plotID).Return(&plotModels.Plot{}, nil).Once()
		mockCameraRepo.On("GetCamerasByPlotID", ctx, plotID).Return(expectedCameras, nil).Once()

		cameras, err := cameraSvc.GetCamerasByPlotID(ctx, plotID)

		assert.NoError(t, err)
		assert.Equal(t, expectedCameras, cameras)
		mockPlotSvc.AssertExpectations(t)
		mockCameraRepo.AssertExpectations(t)
	})

	t.Run("DeleteCamera - Success", func(t *testing.T) {
		cameraID := uuid.New()
		mockCameraRepo.On("DeleteCamera", ctx, cameraID).Return(nil).Once()

		err := cameraSvc.DeleteCamera(ctx, cameraID)

		assert.NoError(t, err)
		mockCameraRepo.AssertExpectations(t)
	})
}
