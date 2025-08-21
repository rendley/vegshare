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


