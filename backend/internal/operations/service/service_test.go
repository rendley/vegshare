package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/rendley/vegshare/backend/internal/catalog/models"
	farmModels "github.com/rendley/vegshare/backend/internal/farm/models"
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

// MockFarmRepository mocks the farm repository
type MockFarmRepository struct {
	mock.Mock
}

// Implement the full farm.Repository interface
func (m *MockFarmRepository) CreateRegion(ctx context.Context, region *farmModels.Region) error {
	return m.Called(ctx, region).Error(0)
}
func (m *MockFarmRepository) GetRegionByID(ctx context.Context, id uuid.UUID) (*farmModels.Region, error) {
	return nil, m.Called(ctx, id).Error(1)
}
func (m *MockFarmRepository) GetAllRegions(ctx context.Context) ([]farmModels.Region, error) {
	return nil, m.Called(ctx).Error(1)
}
func (m *MockFarmRepository) UpdateRegion(ctx context.Context, region *farmModels.Region) error {
	return m.Called(ctx, region).Error(0)
}
func (m *MockFarmRepository) DeleteRegion(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockFarmRepository) CreateLandParcel(ctx context.Context, parcel *farmModels.LandParcel) error {
	return m.Called(ctx, parcel).Error(0)
}
func (m *MockFarmRepository) GetLandParcelByID(ctx context.Context, id uuid.UUID) (*farmModels.LandParcel, error) {
	return nil, m.Called(ctx, id).Error(1)
}
func (m *MockFarmRepository) GetLandParcelsByRegion(ctx context.Context, regionID uuid.UUID) ([]farmModels.LandParcel, error) {
	return nil, m.Called(ctx, regionID).Error(1)
}
func (m *MockFarmRepository) UpdateLandParcel(ctx context.Context, parcel *farmModels.LandParcel) error {
	return m.Called(ctx, parcel).Error(0)
}
func (m *MockFarmRepository) DeleteLandParcel(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockFarmRepository) CreateGreenhouse(ctx context.Context, greenhouse *farmModels.Greenhouse) error {
	return m.Called(ctx, greenhouse).Error(0)
}
func (m *MockFarmRepository) GetGreenhouseByID(ctx context.Context, id uuid.UUID) (*farmModels.Greenhouse, error) {
	return nil, m.Called(ctx, id).Error(1)
}
func (m *MockFarmRepository) GetGreenhousesByLandParcel(ctx context.Context, landParcelID uuid.UUID) ([]farmModels.Greenhouse, error) {
	return nil, m.Called(ctx, landParcelID).Error(1)
}
func (m *MockFarmRepository) UpdateGreenhouse(ctx context.Context, greenhouse *farmModels.Greenhouse) error {
	return m.Called(ctx, greenhouse).Error(0)
}
func (m *MockFarmRepository) DeleteGreenhouse(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockFarmRepository) CreatePlot(ctx context.Context, plot *farmModels.Plot) error {
	return m.Called(ctx, plot).Error(0)
}
func (m *MockFarmRepository) GetPlotsByGreenhouse(ctx context.Context, greenhouseID uuid.UUID) ([]farmModels.Plot, error) {
	return nil, m.Called(ctx, greenhouseID).Error(1)
}
func (m *MockFarmRepository) DeletePlot(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockFarmRepository) GetPlotByID(ctx context.Context, id uuid.UUID) (*farmModels.Plot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*farmModels.Plot), args.Error(1)
}
func (m *MockFarmRepository) UpdatePlot(ctx context.Context, plot *farmModels.Plot) error {
	return m.Called(ctx, plot).Error(0)
}

// MockLeasingRepository mocks the leasing repository
type MockLeasingRepository struct {
	mock.Mock
}

func (m *MockLeasingRepository) CreateLease(ctx context.Context, lease *farmModels.PlotLease) error {
	return m.Called(ctx, lease).Error(0)
}
func (m *MockLeasingRepository) GetLeasesByUserID(ctx context.Context, userID uuid.UUID) ([]farmModels.PlotLease, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]farmModels.PlotLease), args.Error(1)
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
func (m *MockCatalogService) GetAllCrops(ctx context.Context) ([]models.Crop, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Crop), args.Error(1)
}
func (m *MockCatalogService) CreateCrop(ctx context.Context, name, description string, plantingTime, harvestTime int) (*models.Crop, error) {
	args := m.Called(ctx, name, description, plantingTime, harvestTime)
	return args.Get(0).(*models.Crop), args.Error(1)
}

// MockRabbitMQClient mocks the RabbitMQ client
type MockRabbitMQClient struct {
	mock.Mock
}

func (m *MockRabbitMQClient) Publish(queueName, body string) error {
	args := m.Called(queueName, body)
	return args.Error(0)
}
func (m *MockRabbitMQClient) Consume(queueName string) (<-chan amqp.Delivery, error) {
	args := m.Called(queueName)
	return args.Get(0).(<-chan amqp.Delivery), args.Error(1)
}
func (m *MockRabbitMQClient) Close() {
	m.Called()
}

// --- Tests ---

func TestOperationsService(t *testing.T) {
	ctx := context.Background()
	mockOpsRepo := new(MockOperationsRepository)
	mockFarmRepo := new(MockFarmRepository)
	mockLeasingRepo := new(MockLeasingRepository)
	mockCatalogSvc := new(MockCatalogService)
	mockRabbitMQ := new(MockRabbitMQClient)

	opsSvc := NewOperationsService(mockOpsRepo, mockFarmRepo, mockLeasingRepo, mockCatalogSvc, mockRabbitMQ)

	t.Run("PlantCrop", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Arrange
			userID := uuid.New()
			plotID := uuid.New()
			cropID := uuid.New()
			leaseID := uuid.New()
			activeLease := []farmModels.PlotLease{{ID: leaseID, PlotID: plotID, UserID: userID, Status: "active"}}
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
			activeLease := []farmModels.PlotLease{}

			mockLeasingRepo.On("GetLeasesByUserID", ctx, userID).Return(activeLease, nil).Once()

			// Act
			plotCrop, err := opsSvc.PlantCrop(ctx, userID, plotID, cropID)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, plotCrop)
			mockLeasingRepo.AssertExpectations(t)
		})
	})

	t.Run("PerformAction", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Arrange
			plotID := uuid.New()
			action := "water"
			mockRabbitMQ.On("Publish", "actions", mock.AnythingOfType("string")).Return(nil).Once()

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
			mockRabbitMQ.On("Publish", "actions", mock.AnythingOfType("string")).Return(expectedErr).Once()

			// Act
			err := opsSvc.PerformAction(ctx, plotID, action)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, expectedErr, err)
			mockRabbitMQ.AssertExpectations(t)
		})
	})
}