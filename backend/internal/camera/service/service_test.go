package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/camera/models"
	farmModels "github.com/rendley/vegshare/backend/internal/farm/models"
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

// MockFarmRepository is a mock for the farm repository
type MockFarmRepository struct {
	mock.Mock
}

// Implement the full farm.Repository interface needed for these tests
func (m *MockFarmRepository) GetPlotByID(ctx context.Context, id uuid.UUID) (*farmModels.Plot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*farmModels.Plot), args.Error(1)
}

// Dummy implementations for other methods to satisfy the interface
func (m *MockFarmRepository) CreateRegion(ctx context.Context, region *farmModels.Region) error { return nil }
func (m *MockFarmRepository) GetRegionByID(ctx context.Context, id uuid.UUID) (*farmModels.Region, error) { return nil, nil }
func (m *MockFarmRepository) GetAllRegions(ctx context.Context) ([]farmModels.Region, error) { return nil, nil }
func (m *MockFarmRepository) UpdateRegion(ctx context.Context, region *farmModels.Region) error { return nil }
func (m *MockFarmRepository) DeleteRegion(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockFarmRepository) CreateLandParcel(ctx context.Context, parcel *farmModels.LandParcel) error { return nil }
func (m *MockFarmRepository) GetLandParcelByID(ctx context.Context, id uuid.UUID) (*farmModels.LandParcel, error) { return nil, nil }
func (m *MockFarmRepository) GetLandParcelsByRegion(ctx context.Context, regionID uuid.UUID) ([]farmModels.LandParcel, error) { return nil, nil }
func (m *MockFarmRepository) UpdateLandParcel(ctx context.Context, parcel *farmModels.LandParcel) error { return nil }
func (m *MockFarmRepository) DeleteLandParcel(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockFarmRepository) CreateGreenhouse(ctx context.Context, gh *farmModels.Greenhouse) error { return nil }
func (m *MockFarmRepository) GetGreenhouseByID(ctx context.Context, id uuid.UUID) (*farmModels.Greenhouse, error) { return nil, nil }
func (m *MockFarmRepository) GetGreenhousesByLandParcel(ctx context.Context, landParcelID uuid.UUID) ([]farmModels.Greenhouse, error) { return nil, nil }
func (m *MockFarmRepository) UpdateGreenhouse(ctx context.Context, gh *farmModels.Greenhouse) error { return nil }
func (m *MockFarmRepository) DeleteGreenhouse(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockFarmRepository) CreatePlot(ctx context.Context, plot *farmModels.Plot) error { return nil }
func (m *MockFarmRepository) GetPlotsByGreenhouse(ctx context.Context, greenhouseID uuid.UUID) ([]farmModels.Plot, error) { return nil, nil }
func (m *MockFarmRepository) UpdatePlot(ctx context.Context, plot *farmModels.Plot) error { return nil }
func (m *MockFarmRepository) DeletePlot(ctx context.Context, id uuid.UUID) error { return nil }


func TestCameraService(t *testing.T) {
	ctx := context.Background()
	mockCameraRepo := new(MockCameraRepository)
	mockFarmRepo := new(MockFarmRepository)
	cameraSvc := NewService(mockCameraRepo, mockFarmRepo)

	t.Run("CreateCamera - Success", func(t *testing.T) {
		plotID := uuid.New()
		mockFarmRepo.On("GetPlotByID", ctx, plotID).Return(&farmModels.Plot{}, nil).Once()
		mockCameraRepo.On("CreateCamera", ctx, mock.AnythingOfType("*models.Camera")).Return(nil).Once()

		camera, err := cameraSvc.CreateCamera(ctx, "Test Cam", "test_cam", plotID)

		assert.NoError(t, err)
		assert.NotNil(t, camera)
		assert.Equal(t, "Test Cam", camera.Name)
		mockFarmRepo.AssertExpectations(t)
		mockCameraRepo.AssertExpectations(t)
	})

	t.Run("CreateCamera - Plot Not Found", func(t *testing.T) {
		plotID := uuid.New()
		mockFarmRepo.On("GetPlotByID", ctx, plotID).Return(nil, errors.New("not found")).Once()

		camera, err := cameraSvc.CreateCamera(ctx, "Test Cam", "test_cam", plotID)

		assert.Error(t, err)
		assert.Nil(t, camera)
		mockFarmRepo.AssertExpectations(t)
	})

	t.Run("GetCamerasByPlotID - Success", func(t *testing.T) {
		plotID := uuid.New()
		expectedCameras := []models.Camera{{ID: uuid.New()}, {ID: uuid.New()}}
		mockFarmRepo.On("GetPlotByID", ctx, plotID).Return(&farmModels.Plot{}, nil).Once()
		mockCameraRepo.On("GetCamerasByPlotID", ctx, plotID).Return(expectedCameras, nil).Once()

		cameras, err := cameraSvc.GetCamerasByPlotID(ctx, plotID)

		assert.NoError(t, err)
		assert.Equal(t, expectedCameras, cameras)
		mockFarmRepo.AssertExpectations(t)
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