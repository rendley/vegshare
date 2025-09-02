package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rendley/vegshare/backend/internal/catalog/models"
	leasingDomain "github.com/rendley/vegshare/backend/internal/leasing/domain"
	leasingModels "github.com/rendley/vegshare/backend/internal/leasing/models"
	plotModels "github.com/rendley/vegshare/backend/internal/plot/models"
	"github.com/rendley/vegshare/backend/pkg/config"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

// MockOperationsRepository mocks the operations repository
type MockOperationsRepository struct {
	mock.Mock
}

func (m *MockOperationsRepository) CreatePlotCrop(ctx context.Context, plotCrop *models.PlotCrop) error {
	args := m.Called(ctx, plotCrop)
	return args.Error(0)
}
func (m *MockOperationsRepository) GetPlotCrops(ctx context.Context, plotID uuid.UUID) ([]models.PlotCrop, error) {
	args := m.Called(ctx, plotID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PlotCrop), args.Error(1)
}
func (m *MockOperationsRepository) DeletePlotCrop(ctx context.Context, plantingID uuid.UUID) error {
	args := m.Called(ctx, plantingID)
	return args.Error(0)
}

// MockPlotService mocks the plot service
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
func (m *MockPlotService) CreatePlot(ctx context.Context, name, size string, greenhouseID uuid.UUID) (*plotModels.Plot, error) {
	return nil, nil
}
func (m *MockPlotService) GetPlotsByGreenhouse(ctx context.Context, greenhouseID uuid.UUID) ([]plotModels.Plot, error) {
	return nil, nil
}
func (m *MockPlotService) DeletePlot(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockPlotService) GetLeasableUnit(ctx context.Context, unitID uuid.UUID) (leasingDomain.LeasableUnit, error) {
	return nil, nil
}
func (m *MockPlotService) UpdateUnitStatus(ctx context.Context, unitID uuid.UUID, status string) error {
	return nil
}
func (m *MockPlotService) WithTx(tx *sqlx.Tx) leasingDomain.UnitManager { return m }

// MockLeasingRepository mocks the leasing repository
type MockLeasingRepository struct {
	mock.Mock
}

func (m *MockLeasingRepository) CreateLease(ctx context.Context, lease *leasingModels.Lease) error {
	return m.Called(ctx, lease).Error(0)
}
func (m *MockLeasingRepository) GetLeasesByUserID(ctx context.Context, userID uuid.UUID) ([]leasingModels.Lease, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]leasingModels.Lease), args.Error(1)
}

// MockCatalogService mocks the catalog service
type MockCatalogService struct {
	mock.Mock
}

func (m *MockCatalogService) GetCropByID(ctx context.Context, id uuid.UUID) (*models.Crop, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Crop), args.Error(1)
}
func (m *MockCatalogService) GetAllCrops(ctx context.Context) ([]models.Crop, error) { return nil, nil }
func (m *MockCatalogService) CreateCrop(ctx context.Context, name, description string, plantingTime, harvestTime int) (*models.Crop, error) {
	return nil, nil
}

// MockRabbitMQClient mocks the RabbitMQ client
type MockRabbitMQClient struct {
	mock.Mock
}

func (m *MockRabbitMQClient) Publish(queueName, body string) error {
	return m.Called(queueName, body).Error(0)
}
func (m *MockRabbitMQClient) Consume(queueName string) (<-chan amqp.Delivery, error) { return nil, nil }
func (m *MockRabbitMQClient) Close()                                                 { m.Called() }

// --- Tests ---

func TestOperationsService(t *testing.T) {
	ctx := context.Background()
	mockOpsRepo := new(MockOperationsRepository)
	mockPlotSvc := new(MockPlotService)
	mockLeasingRepo := new(MockLeasingRepository)
	mockCatalogSvc := new(MockCatalogService)
	mockRabbitMQ := new(MockRabbitMQClient)
	cfg := &config.Config{
		RabbitMQ: config.RabbitMQConfig{
			Queues: map[string]string{"actions": "actions_queue_test"},
		},
	}

	opsSvc := NewOperationsService(mockOpsRepo, mockPlotSvc, mockLeasingRepo, mockCatalogSvc, mockRabbitMQ, cfg)

	t.Run("PlantCrop", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Arrange
			userID := uuid.New()
			plotID := uuid.New()
			cropID := uuid.New()
			leaseID := uuid.New()
			activeLease := []leasingModels.Lease{{ID: leaseID, UnitID: plotID, UserID: userID, Status: "active", UnitType: leasingModels.UnitTypePlot}}
			crop := &models.Crop{ID: cropID}

			mockLeasingRepo.On("GetLeasesByUserID", ctx, userID).Return(activeLease, nil).Once()
			mockCatalogSvc.On("GetCropByID", ctx, cropID).Return(crop, nil).Once()
			mockOpsRepo.On("CreatePlotCrop", ctx, mock.AnythingOfType("*models.PlotCrop")).Return(nil).Once()

			// Act
			plotCrop, err := opsSvc.PlantCrop(ctx, userID, plotID, cropID)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, plotCrop)
			assert.Equal(t, cropID, plotCrop.CropID)
			mockLeasingRepo.AssertExpectations(t)
			mockCatalogSvc.AssertExpectations(t)
			mockOpsRepo.AssertExpectations(t)
		})

		t.Run("No active lease", func(t *testing.T) {
			// Arrange
			userID := uuid.New()
			plotID := uuid.New()
			cropID := uuid.New()
			activeLease := []leasingModels.Lease{}

			mockLeasingRepo.On("GetLeasesByUserID", ctx, userID).Return(activeLease, nil).Once()

			// Act
			plotCrop, err := opsSvc.PlantCrop(ctx, userID, plotID, cropID)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, plotCrop)
			mockLeasingRepo.AssertExpectations(t)
		})
	})

	// ... (остальные тесты для PerformAction без изменений) ...
	t.Run("PerformAction", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Arrange
			plotID := uuid.New()
			action := "water"
			mockRabbitMQ.On("Publish", cfg.RabbitMQ.Queues["actions"], mock.AnythingOfType("string")).Return(nil).Once()

			// Act
			err := opsSvc.PerformAction(ctx, plotID, action)

			// Assert
			assert.NoError(t, err)
			mockRabbitMQ.AssertExpectations(t)
		})

		t.Run("Error on publish", func(t *testing.T) {
			// Arrange
			plotID := uuid.New()
			action := "water"
			expectedErr := errors.New("publish error")
			mockRabbitMQ.On("Publish", cfg.RabbitMQ.Queues["actions"], mock.AnythingOfType("string")).Return(expectedErr).Once()

			// Act
			err := opsSvc.PerformAction(ctx, plotID, action)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, expectedErr, err)
			mockRabbitMQ.AssertExpectations(t)
		})
	})
}
