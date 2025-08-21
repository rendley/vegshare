package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	farmModels "github.com/rendley/vegshare/backend/internal/farm/models"
	"github.com/rendley/vegshare/backend/internal/plot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockPlotRepository struct {
	mock.Mock
}

func (m *MockPlotRepository) CreatePlot(ctx context.Context, plot *models.Plot) error {
	args := m.Called(ctx, plot)
	return args.Error(0)
}

func (m *MockPlotRepository) GetPlotByID(ctx context.Context, id uuid.UUID) (*models.Plot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Plot), args.Error(1)
}

func (m *MockPlotRepository) GetPlotsByGreenhouse(ctx context.Context, greenhouseID uuid.UUID) ([]models.Plot, error) {
	args := m.Called(ctx, greenhouseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Plot), args.Error(1)
}

func (m *MockPlotRepository) UpdatePlot(ctx context.Context, plot *models.Plot) error {
	args := m.Called(ctx, plot)
	return args.Error(0)
}

func (m *MockPlotRepository) DeletePlot(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockFarmService struct {
	mock.Mock
}

func (m *MockFarmService) GetGreenhouseByID(ctx context.Context, id uuid.UUID) (*farmModels.Greenhouse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*farmModels.Greenhouse), args.Error(1)
}

// Dummy implementations for other farm service methods to satisfy the interface
func (m *MockFarmService) CreateRegion(ctx context.Context, name string) (*farmModels.Region, error) { return nil, nil }
func (m *MockFarmService) GetRegionByID(ctx context.Context, id uuid.UUID) (*farmModels.Region, error) { return nil, nil }
func (m *MockFarmService) GetAllRegions(ctx context.Context) ([]farmModels.Region, error) { return nil, nil }
func (m *MockFarmService) UpdateRegion(ctx context.Context, id uuid.UUID, name string) (*farmModels.Region, error) { return nil, nil }
func (m *MockFarmService) DeleteRegion(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockFarmService) CreateLandParcel(ctx context.Context, name string, regionID uuid.UUID) (*farmModels.LandParcel, error) { return nil, nil }
func (m *MockFarmService) GetLandParcelByID(ctx context.Context, id uuid.UUID) (*farmModels.LandParcel, error) { return nil, nil }
func (m *MockFarmService) GetLandParcelsByRegion(ctx context.Context, regionID uuid.UUID) ([]farmModels.LandParcel, error) { return nil, nil }
func (m *MockFarmService) UpdateLandParcel(ctx context.Context, id uuid.UUID, name string) (*farmModels.LandParcel, error) { return nil, nil }
func (m *MockFarmService) DeleteLandParcel(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockFarmService) CreateGreenhouse(ctx context.Context, name, typeName string, landParcelID uuid.UUID) (*farmModels.Greenhouse, error) { return nil, nil }
func (m *MockFarmService) GetGreenhousesByLandParcel(ctx context.Context, landParcelID uuid.UUID) ([]farmModels.Greenhouse, error) { return nil, nil }
func (m *MockFarmService) UpdateGreenhouse(ctx context.Context, id uuid.UUID, name, typeName string) (*farmModels.Greenhouse, error) { return nil, nil }
func (m *MockFarmService) DeleteGreenhouse(ctx context.Context, id uuid.UUID) error { return nil }


func TestPlotService(t *testing.T) {
	ctx := context.Background()
	mockPlotRepo := new(MockPlotRepository)
	mockFarmSvc := new(MockFarmService)
	plotSvc := NewService(mockPlotRepo, mockFarmSvc)

	t.Run("CreatePlot - Success", func(t *testing.T) {
		greenhouseID := uuid.New()
		plotName := "Test Plot"
		plotSize := "2x2"
		mockFarmSvc.On("GetGreenhouseByID", ctx, greenhouseID).Return(&farmModels.Greenhouse{}, nil).Once()
		mockPlotRepo.On("CreatePlot", ctx, mock.AnythingOfType("*models.Plot")).Return(nil).Once()

		plot, err := plotSvc.CreatePlot(ctx, plotName, plotSize, greenhouseID)

		assert.NoError(t, err)
		assert.NotNil(t, plot)
		assert.Equal(t, plotName, plot.Name)
		assert.Equal(t, plotSize, plot.Size)
		assert.Equal(t, "available", plot.Status)
		assert.Equal(t, greenhouseID, plot.GreenhouseID)
		mockPlotRepo.AssertExpectations(t)
		mockFarmSvc.AssertExpectations(t)
	})

	t.Run("GetPlotByID - Success", func(t *testing.T) {
		plotID := uuid.New()
		expectedPlot := &models.Plot{ID: plotID, Name: "Test Plot"}
		mockPlotRepo.On("GetPlotByID", ctx, plotID).Return(expectedPlot, nil).Once()

		plot, err := plotSvc.GetPlotByID(ctx, plotID)

		assert.NoError(t, err)
		assert.Equal(t, expectedPlot, plot)
		mockPlotRepo.AssertExpectations(t)
	})

	t.Run("UpdatePlot - Success", func(t *testing.T) {
		plotID := uuid.New()
		originalPlot := &models.Plot{ID: plotID, Name: "Original Name", Size: "1x1", Status: "available"}
		mockPlotRepo.On("GetPlotByID", ctx, plotID).Return(originalPlot, nil).Once()
		mockPlotRepo.On("UpdatePlot", ctx, mock.AnythingOfType("*models.Plot")).Return(nil).Once()

		updatedPlot, err := plotSvc.UpdatePlot(ctx, plotID, "Updated Name", "2x2", "rented")

		assert.NoError(t, err)
		assert.NotNil(t, updatedPlot)
		assert.Equal(t, "Updated Name", updatedPlot.Name)
		assert.Equal(t, "2x2", updatedPlot.Size)
		assert.Equal(t, "rented", updatedPlot.Status)
		mockPlotRepo.AssertExpectations(t)
	})

	t.Run("DeletePlot - Success", func(t *testing.T) {
		plotID := uuid.New()
		mockPlotRepo.On("DeletePlot", ctx, plotID).Return(nil).Once()

		err := plotSvc.DeletePlot(ctx, plotID)

		assert.NoError(t, err)
		mockPlotRepo.AssertExpectations(t)
	})

	t.Run("GetPlotsByGreenhouse - Success", func(t *testing.T) {
		greenhouseID := uuid.New()
		expectedPlots := []models.Plot{{ID: uuid.New()}, {ID: uuid.New()}}
		mockPlotRepo.On("GetPlotsByGreenhouse", ctx, greenhouseID).Return(expectedPlots, nil).Once()

		plots, err := plotSvc.GetPlotsByGreenhouse(ctx, greenhouseID)

		assert.NoError(t, err)
		assert.Equal(t, expectedPlots, plots)
		mockPlotRepo.AssertExpectations(t)
	})
}
