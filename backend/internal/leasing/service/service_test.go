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

// --- Mocks ---

// MockLeasingRepository is a mock for the leasing repository
type MockLeasingRepository struct {
	mock.Mock
}

func (m *MockLeasingRepository) CreateLease(ctx context.Context, lease *models.PlotLease) error {
	args := m.Called(ctx, lease)
	return args.Error(0)
}

func (m *MockLeasingRepository) GetLeasesByUserID(ctx context.Context, userID uuid.UUID) ([]models.PlotLease, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PlotLease), args.Error(1)
}

// MockFarmRepository is a mock for the farm repository
type MockFarmRepository struct {
	mock.Mock
}

// Implement the full farm.Repository interface
func (m *MockFarmRepository) CreateRegion(ctx context.Context, region *models.Region) error {
	args := m.Called(ctx, region)
	return args.Error(0)
}
func (m *MockFarmRepository) GetRegionByID(ctx context.Context, id uuid.UUID) (*models.Region, error) {
	args := m.Called(ctx, id)
	return nil, args.Error(1)
}
func (m *MockFarmRepository) GetAllRegions(ctx context.Context) ([]models.Region, error) {
	args := m.Called(ctx)
	return nil, args.Error(1)
}
func (m *MockFarmRepository) UpdateRegion(ctx context.Context, region *models.Region) error {
	args := m.Called(ctx, region)
	return args.Error(0)
}
func (m *MockFarmRepository) DeleteRegion(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockFarmRepository) CreateLandParcel(ctx context.Context, parcel *models.LandParcel) error {
	args := m.Called(ctx, parcel)
	return args.Error(0)
}
func (m *MockFarmRepository) GetLandParcelByID(ctx context.Context, id uuid.UUID) (*models.LandParcel, error) {
	args := m.Called(ctx, id)
	return nil, args.Error(1)
}
func (m *MockFarmRepository) GetLandParcelsByRegion(ctx context.Context, regionID uuid.UUID) ([]models.LandParcel, error) {
	args := m.Called(ctx, regionID)
	return nil, args.Error(1)
}
func (m *MockFarmRepository) UpdateLandParcel(ctx context.Context, parcel *models.LandParcel) error {
	args := m.Called(ctx, parcel)
	return args.Error(0)
}
func (m *MockFarmRepository) DeleteLandParcel(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockFarmRepository) CreateGreenhouse(ctx context.Context, greenhouse *models.Greenhouse) error {
	args := m.Called(ctx, greenhouse)
	return args.Error(0)
}
func (m *MockFarmRepository) GetGreenhouseByID(ctx context.Context, id uuid.UUID) (*models.Greenhouse, error) {
	args := m.Called(ctx, id)
	return nil, args.Error(1)
}
func (m *MockFarmRepository) GetGreenhousesByLandParcel(ctx context.Context, landParcelID uuid.UUID) ([]models.Greenhouse, error) {
	args := m.Called(ctx, landParcelID)
	return nil, args.Error(1)
}
func (m *MockFarmRepository) UpdateGreenhouse(ctx context.Context, greenhouse *models.Greenhouse) error {
	args := m.Called(ctx, greenhouse)
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
func (m *MockFarmRepository) GetPlotsByGreenhouse(ctx context.Context, greenhouseID uuid.UUID) ([]models.Plot, error) {
	args := m.Called(ctx, greenhouseID)
	return nil, args.Error(1)
}
func (m *MockFarmRepository) DeletePlot(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFarmRepository) GetPlotByID(ctx context.Context, id uuid.UUID) (*models.Plot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Plot), args.Error(1)
}

func (m *MockFarmRepository) UpdatePlot(ctx context.Context, plot *models.Plot) error {
	args := m.Called(ctx, plot)
	return args.Error(0)
}

// --- Tests ---

func TestLeasingService(t *testing.T) {
	ctx := context.Background()
	mockLeasingRepo := new(MockLeasingRepository)
	mockFarmRepo := new(MockFarmRepository)

	leasingSvc := NewLeasingService(mockLeasingRepo, mockFarmRepo)

	t.Run("LeasePlot", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Arrange
			userID := uuid.New()
			plotID := uuid.New()
			availablePlot := &models.Plot{ID: plotID, Status: "available"}

			mockFarmRepo.On("GetPlotByID", ctx, plotID).Return(availablePlot, nil).Once()
			mockLeasingRepo.On("CreateLease", ctx, mock.AnythingOfType("*models.PlotLease")).Return(nil).Once()
			mockFarmRepo.On("UpdatePlot", ctx, mock.AnythingOfType("*models.Plot")).Return(nil).Once()

			// Act
			lease, err := leasingSvc.LeasePlot(ctx, userID, plotID)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, lease)
			assert.Equal(t, "active", lease.Status)
			assert.Equal(t, plotID, lease.PlotID)
			assert.Equal(t, userID, lease.UserID)
			mockLeasingRepo.AssertExpectations(t)
			mockFarmRepo.AssertExpectations(t)
		})

		t.Run("Plot not found", func(t *testing.T) {
			// Arrange
			userID := uuid.New()
			plotID := uuid.New()
			mockFarmRepo.On("GetPlotByID", ctx, plotID).Return(nil, errors.New("not found")).Once()

			// Act
			lease, err := leasingSvc.LeasePlot(ctx, userID, plotID)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, lease)
			mockFarmRepo.AssertExpectations(t)
		})

		t.Run("Plot not available", func(t *testing.T) {
			// Arrange
			userID := uuid.New()
			plotID := uuid.New()
			rentedPlot := &models.Plot{ID: plotID, Status: "rented"}

			mockFarmRepo.On("GetPlotByID", ctx, plotID).Return(rentedPlot, nil).Once()

			// Act
			lease, err := leasingSvc.LeasePlot(ctx, userID, plotID)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, lease)
			mockFarmRepo.AssertExpectations(t)
		})
	})

	t.Run("GetMyLeases", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Arrange
			userID := uuid.New()
			expectedLeases := []models.PlotLease{{ID: uuid.New(), UserID: userID}}
			mockLeasingRepo.On("GetLeasesByUserID", ctx, userID).Return(expectedLeases, nil).Once()

			// Act
			leases, err := leasingSvc.GetMyLeases(ctx, userID)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expectedLeases, leases)
			mockLeasingRepo.AssertExpectations(t)
		})

		t.Run("Error", func(t *testing.T) {
			// Arrange
			userID := uuid.New()
			mockLeasingRepo.On("GetLeasesByUserID", ctx, userID).Return(nil, errors.New("db error")).Once()

			// Act
			leases, err := leasingSvc.GetMyLeases(ctx, userID)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, leases)
			mockLeasingRepo.AssertExpectations(t)
		})
	})
}