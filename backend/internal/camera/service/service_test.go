package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/camera/models"
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

func (m *MockCameraRepository) GetCamerasByUnitID(ctx context.Context, unitID uuid.UUID, unitType string) ([]models.Camera, error) {
	args := m.Called(ctx, unitID, unitType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Camera), args.Error(1)
}

func (m *MockCameraRepository) DeleteCamera(ctx context.Context, cameraID uuid.UUID) error {
	args := m.Called(ctx, cameraID)
	return args.Error(0)
}

func (m *MockCameraRepository) GetCameraByID(ctx context.Context, cameraID uuid.UUID) (*models.Camera, error) {
	args := m.Called(ctx, cameraID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Camera), args.Error(1)
}

// --- Tests ---

func TestCameraService(t *testing.T) {
	ctx := context.Background()
	mockCameraRepo := new(MockCameraRepository)
	cameraSvc := NewService(mockCameraRepo)

	t.Run("CreateCamera - Success", func(t *testing.T) {
		unitID := uuid.New()
		unitType := "plot"
		mockCameraRepo.On("CreateCamera", ctx, mock.AnythingOfType("*models.Camera")).Return(nil).Once()

		camera, err := cameraSvc.CreateCamera(ctx, "Test Cam", "test_cam", unitID, unitType)

		assert.NoError(t, err)
		assert.NotNil(t, camera)
		assert.Equal(t, "Test Cam", camera.Name)
		assert.Equal(t, unitID, camera.UnitID)
		assert.Equal(t, unitType, camera.UnitType)
		mockCameraRepo.AssertExpectations(t)
	})

	t.Run("GetCamerasByUnitID - Success", func(t *testing.T) {
		unitID := uuid.New()
		unitType := "plot"
		expectedCameras := []models.Camera{{ID: uuid.New()}, {ID: uuid.New()}}
		mockCameraRepo.On("GetCamerasByUnitID", ctx, unitID, unitType).Return(expectedCameras, nil).Once()

		cameras, err := cameraSvc.GetCamerasByUnitID(ctx, unitID, unitType)

		assert.NoError(t, err)
		assert.Equal(t, expectedCameras, cameras)
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