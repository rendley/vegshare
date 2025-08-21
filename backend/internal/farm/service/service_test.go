package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/farm/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock Repository ---

type MockFarmRepository struct {
	mock.Mock
}

// Implement all methods of the repository.Repository interface



func (m *MockFarmRepository) CreateRegion(ctx context.Context, region *models.Region) error {
	args := m.Called(ctx, region)
	return args.Error(0)
}

func (m *MockFarmRepository) GetRegionByID(ctx context.Context, id uuid.UUID) (*models.Region, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Region), args.Error(1)
}

func (m *MockFarmRepository) GetAllRegions(ctx context.Context) ([]models.Region, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Region), args.Error(1)
}

func (m *MockFarmRepository) UpdateRegion(ctx context.Context, region *models.Region) error {
	args := m.Called(ctx, region)
	return args.Error(0)
}

func (m *MockFarmRepository) DeleteRegion(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ... Implement mock methods for LandParcel, Greenhouse, and Plot in a similar fashion ...

func (m *MockFarmRepository) CreateLandParcel(ctx context.Context, parcel *models.LandParcel) error {
	args := m.Called(ctx, parcel)
	return args.Error(0)
}
func (m *MockFarmRepository) GetLandParcelByID(ctx context.Context, id uuid.UUID) (*models.LandParcel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LandParcel), args.Error(1)
}
func (m *MockFarmRepository) GetLandParcelsByRegion(ctx context.Context, regionID uuid.UUID) ([]models.LandParcel, error) {
	args := m.Called(ctx, regionID)
	return args.Get(0).([]models.LandParcel), args.Error(1)
}
func (m *MockFarmRepository) UpdateLandParcel(ctx context.Context, parcel *models.LandParcel) error {
	args := m.Called(ctx, parcel)
	return args.Error(0)
}
func (m *MockFarmRepository) DeleteLandParcel(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFarmRepository) CreateGreenhouse(ctx context.Context, gh *models.Greenhouse) error {
	args := m.Called(ctx, gh)
	return args.Error(0)
}
func (m *MockFarmRepository) GetGreenhouseByID(ctx context.Context, id uuid.UUID) (*models.Greenhouse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Greenhouse), args.Error(1)
}
func (m *MockFarmRepository) GetGreenhousesByLandParcel(ctx context.Context, landParcelID uuid.UUID) ([]models.Greenhouse, error) {
	args := m.Called(ctx, landParcelID)
	return args.Get(0).([]models.Greenhouse), args.Error(1)
}
func (m *MockFarmRepository) UpdateGreenhouse(ctx context.Context, gh *models.Greenhouse) error {
	args := m.Called(ctx, gh)
	return args.Error(0)
}
func (m *MockFarmRepository) DeleteGreenhouse(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFarmRepository) CreatePlot(ctx context.Context, plot *models.Plot) error {
	args := m.Called(ctx, plot)
	return args.Error(0)
}
func (m *MockFarmRepository) GetPlotByID(ctx context.Context, id uuid.UUID) (*models.Plot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Plot), args.Error(1)
}
func (m *MockFarmRepository) GetPlotsByGreenhouse(ctx context.Context, greenhouseID uuid.UUID) ([]models.Plot, error) {
	args := m.Called(ctx, greenhouseID)
	return args.Get(0).([]models.Plot), args.Error(1)
}
func (m *MockFarmRepository) UpdatePlot(ctx context.Context, plot *models.Plot) error {
	args := m.Called(ctx, plot)
	return args.Error(0)
}
func (m *MockFarmRepository) DeletePlot(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --- Tests ---

func TestFarmService_Regions(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockFarmRepository)
	farmSvc := NewFarmService(mockRepo)

	t.Run("CreateRegion - Success", func(t *testing.T) {
		// Arrange
		mockRepo.On("CreateRegion", ctx, mock.AnythingOfType("*models.Region")).Return(nil).Once()

		// Act
		region, err := farmSvc.CreateRegion(ctx, "Test Region")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, region)
		assert.Equal(t, "Test Region", region.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("CreateRegion - Error", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("db error")
		mockRepo.On("CreateRegion", ctx, mock.AnythingOfType("*models.Region")).Return(expectedErr).Once()

		// Act
		region, err := farmSvc.CreateRegion(ctx, "Test Region")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, region)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetRegionByID - Success", func(t *testing.T) {
		// Arrange
		regionID := uuid.New()
		expectedRegion := &models.Region{ID: regionID, Name: "Test"}
		mockRepo.On("GetRegionByID", ctx, regionID).Return(expectedRegion, nil).Once()

		// Act
		region, err := farmSvc.GetRegionByID(ctx, regionID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedRegion, region)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetRegionByID - Error", func(t *testing.T) {
		// Arrange
		regionID := uuid.New()
		expectedErr := errors.New("not found")
		mockRepo.On("GetRegionByID", ctx, regionID).Return(nil, expectedErr).Once()

		// Act
		region, err := farmSvc.GetRegionByID(ctx, regionID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, region)
		mockRepo.AssertExpectations(t)
	})
}

// We can add more tests for other entities (LandParcel, Greenhouse, Plot, Crop) in a similar manner.

func TestFarmService_Plots(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockFarmRepository)
	farmSvc := NewFarmService(mockRepo)

	t.Run("CreatePlot - Success", func(t *testing.T) {
		// Arrange
		greenhouseID := uuid.New()
		plotName := "Test Plot"
		plotSize := "2x2"
		mockRepo.On("GetGreenhouseByID", ctx, greenhouseID).Return(&models.Greenhouse{}, nil).Once()
		mockRepo.On("CreatePlot", ctx, mock.AnythingOfType("*models.Plot")).Return(nil).Once()

		// Act
		plot, err := farmSvc.CreatePlot(ctx, plotName, plotSize, greenhouseID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, plot)
		assert.Equal(t, plotName, plot.Name)
		assert.Equal(t, plotSize, plot.Size)
		assert.Equal(t, "available", plot.Status)
		assert.Equal(t, greenhouseID, plot.GreenhouseID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetPlotByID - Success", func(t *testing.T) {
		// Arrange
		plotID := uuid.New()
		expectedPlot := &models.Plot{ID: plotID, Name: "Test Plot"}
		mockRepo.On("GetPlotByID", ctx, plotID).Return(expectedPlot, nil).Once()

		// Act
		plot, err := farmSvc.GetPlotByID(ctx, plotID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedPlot, plot)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdatePlot - Success", func(t *testing.T) {
		// Arrange
		plotID := uuid.New()
		originalPlot := &models.Plot{ID: plotID, Name: "Original Name", Size: "1x1", Status: "available"}
		mockRepo.On("GetPlotByID", ctx, plotID).Return(originalPlot, nil).Once()
		mockRepo.On("UpdatePlot", ctx, mock.AnythingOfType("*models.Plot")).Return(nil).Once()

		// Act
		updatedPlot, err := farmSvc.UpdatePlot(ctx, plotID, "Updated Name", "2x2", "rented")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedPlot)
		assert.Equal(t, "Updated Name", updatedPlot.Name)
		assert.Equal(t, "2x2", updatedPlot.Size)
		assert.Equal(t, "rented", updatedPlot.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DeletePlot - Success", func(t *testing.T) {
		// Arrange
		plotID := uuid.New()
		mockRepo.On("DeletePlot", ctx, plotID).Return(nil).Once()

		// Act
		err := farmSvc.DeletePlot(ctx, plotID)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetPlotsByGreenhouse - Success", func(t *testing.T) {
		// Arrange
		greenhouseID := uuid.New()
		expectedPlots := []models.Plot{
			{ID: uuid.New(), Name: "Plot 1", GreenhouseID: greenhouseID},
			{ID: uuid.New(), Name: "Plot 2", GreenhouseID: greenhouseID},
		}
		mockRepo.On("GetPlotsByGreenhouse", ctx, greenhouseID).Return(expectedPlots, nil).Once()

		// Act
		plots, err := farmSvc.GetPlotsByGreenhouse(ctx, greenhouseID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedPlots, plots)
		mockRepo.AssertExpectations(t)
	})
}
