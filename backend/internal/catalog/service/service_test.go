package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/catalog/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock Repository ---

type MockCatalogRepository struct {
	mock.Mock
}

func (m *MockCatalogRepository) GetAllCrops(ctx context.Context) ([]models.Crop, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Crop), args.Error(1)
}

func (m *MockCatalogRepository) GetCropByID(ctx context.Context, id uuid.UUID) (*models.Crop, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Crop), args.Error(1)
}

func (m *MockCatalogRepository) CreateCrop(ctx context.Context, crop *models.Crop) error {
	args := m.Called(ctx, crop)
	return args.Error(0)
}

// --- Tests ---

func TestCatalogService_Crops(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockCatalogRepository)
	catalogSvc := NewService(mockRepo)

	t.Run("CreateCrop", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Arrange
			mockRepo.On("CreateCrop", ctx, mock.AnythingOfType("*models.Crop")).Return(nil).Once()

			// Act
			crop, err := catalogSvc.CreateCrop(ctx, "Tomato", "Red and juicy", 30, 90)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, crop)
			assert.Equal(t, "Tomato", crop.Name)
			mockRepo.AssertExpectations(t)
		})

		t.Run("Error", func(t *testing.T) {
			// Arrange
			expectedErr := errors.New("db error")
			mockRepo.On("CreateCrop", ctx, mock.AnythingOfType("*models.Crop")).Return(expectedErr).Once()

			// Act
			crop, err := catalogSvc.CreateCrop(ctx, "Tomato", "Red and juicy", 30, 90)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, crop)
			assert.Equal(t, expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	})

	t.Run("GetCropByID", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Arrange
			cropID := uuid.New()
			expectedCrop := &models.Crop{ID: cropID, Name: "Carrot"}
			mockRepo.On("GetCropByID", ctx, cropID).Return(expectedCrop, nil).Once()

			// Act
			crop, err := catalogSvc.GetCropByID(ctx, cropID)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expectedCrop, crop)
			mockRepo.AssertExpectations(t)
		})

		t.Run("Error", func(t *testing.T) {
			// Arrange
			cropID := uuid.New()
			expectedErr := errors.New("not found")
			mockRepo.On("GetCropByID", ctx, cropID).Return(nil, expectedErr).Once()

			// Act
			crop, err := catalogSvc.GetCropByID(ctx, cropID)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, crop)
			mockRepo.AssertExpectations(t)
		})
	})

	t.Run("GetAllCrops", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Arrange
			expectedCrops := []models.Crop{{ID: uuid.New(), Name: "Lettuce"}}
			mockRepo.On("GetAllCrops", ctx).Return(expectedCrops, nil).Once()

			// Act
			crops, err := catalogSvc.GetAllCrops(ctx)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expectedCrops, crops)
			mockRepo.AssertExpectations(t)
		})

		t.Run("Error", func(t *testing.T) {
			// Arrange
			expectedErr := errors.New("db error")
			mockRepo.On("GetAllCrops", ctx).Return(nil, expectedErr).Once()

			// Act
			crops, err := catalogSvc.GetAllCrops(ctx)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, crops)
			mockRepo.AssertExpectations(t)
		})
	})
}