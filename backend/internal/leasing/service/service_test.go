package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rendley/vegshare/backend/internal/leasing/models"
	plotModels "github.com/rendley/vegshare/backend/internal/plot/models"
	plotService "github.com/rendley/vegshare/backend/internal/plot/service"
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

func (m *MockPlotService) UpdatePlot(ctx context.Context, id uuid.UUID, name, size, status string) (*plotModels.Plot, error) {
	args := m.Called(ctx, id, name, size, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*plotModels.Plot), args.Error(1)
}

func (m *MockPlotService) WithTx(tx *sqlx.Tx) plotService.Service {
	// In tests, we can often just return the same mock
	// if we don't need to test the transactional behavior itself.
	return m
}

// Dummy implementations for other plot service methods
func (m *MockPlotService) CreatePlot(ctx context.Context, name, size string, greenhouseID uuid.UUID) (*plotModels.Plot, error) { return nil, nil }
func (m *MockPlotService) GetPlotsByGreenhouse(ctx context.Context, greenhouseID uuid.UUID) ([]plotModels.Plot, error) { return nil, nil }
func (m *MockPlotService) DeletePlot(ctx context.Context, id uuid.UUID) error { return nil }


// --- Tests ---

func TestLeasingService(t *testing.T) {
	ctx := context.Background()
	mockLeasingRepo := new(MockLeasingRepository)
	mockPlotSvc := new(MockPlotService)

	// We pass nil for the DB connection, as unit tests for business logic
	// shouldn't rely on a real database connection.
	leasingSvc := NewLeasingService(nil, mockLeasingRepo, mockPlotSvc)

	t.Run("LeasePlot", func(t *testing.T) {
		t.Run("Plot not found", func(t *testing.T) {
			// Arrange
			userID := uuid.New()
			plotID := uuid.New()
			mockPlotSvc.On("GetPlotByID", ctx, plotID).Return(nil, errors.New("not found")).Once()

			// Act
			lease, err := leasingSvc.LeasePlot(ctx, userID, plotID)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, lease)
			mockPlotSvc.AssertExpectations(t)
		})

		t.Run("Plot not available", func(t *testing.T) {
			// Arrange
			userID := uuid.New()
			plotID := uuid.New()
			rentedPlot := &plotModels.Plot{ID: plotID, Status: "rented"}

			mockPlotSvc.On("GetPlotByID", ctx, plotID).Return(rentedPlot, nil).Once()

			// Act
			lease, err := leasingSvc.LeasePlot(ctx, userID, plotID)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, lease)
			mockPlotSvc.AssertExpectations(t)
		})

		// Note: A successful transaction test is an integration test, not a unit test.
		// The unit test for the success case would fail on `db.BeginTxx` because db is nil.
		// We are only testing the business logic branches here.
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
