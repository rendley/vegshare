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

func (m *MockCatalogRepository) GetItemsByType(ctx context.Context, itemType string) ([]models.CatalogItem, error) {
	args := m.Called(ctx, itemType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.CatalogItem), args.Error(1)
}

func (m *MockCatalogRepository) GetItemByID(ctx context.Context, id uuid.UUID) (*models.CatalogItem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CatalogItem), args.Error(1)
}

func (m *MockCatalogRepository) CreateItem(ctx context.Context, item *models.CatalogItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

// --- Tests ---

func TestCatalogService_Items(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockCatalogRepository)
	catalogSvc := NewService(mockRepo)

	t.Run("CreateItem", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Arrange
			mockRepo.On("CreateItem", ctx, mock.AnythingOfType("*models.CatalogItem")).Return(nil).Once()
			attributes := models.JSONB{"planting_time": 30, "harvest_time": 90}

			// Act
			item, err := catalogSvc.CreateItem(ctx, "crop", "Tomato", "Red and juicy", attributes)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, item)
			assert.Equal(t, "crop", item.ItemType)
			assert.Equal(t, "Tomato", item.Name)
			assert.Equal(t, attributes, item.Attributes)
			mockRepo.AssertExpectations(t)
		})

		t.Run("Error on DB", func(t *testing.T) {
			// Arrange
			expectedErr := errors.New("db error")
			mockRepo.On("CreateItem", ctx, mock.AnythingOfType("*models.CatalogItem")).Return(expectedErr).Once()

			// Act
			item, err := catalogSvc.CreateItem(ctx, "crop", "Tomato", "Red and juicy", nil)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, item)
			assert.Equal(t, expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	})

	t.Run("GetItemByID", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Arrange
			itemID := uuid.New()
			expectedItem := &models.CatalogItem{ID: itemID, Name: "Carrot", ItemType: "crop"}
			mockRepo.On("GetItemByID", ctx, itemID).Return(expectedItem, nil).Once()

			// Act
			item, err := catalogSvc.GetItemByID(ctx, itemID)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expectedItem, item)
			mockRepo.AssertExpectations(t)
		})
	})

	t.Run("GetItems", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Arrange
			expectedItems := []models.CatalogItem{{ID: uuid.New(), Name: "Lettuce", ItemType: "crop"}}
			mockRepo.On("GetItemsByType", ctx, "crop").Return(expectedItems, nil).Once()

			// Act
			items, err := catalogSvc.GetItems(ctx, "crop")

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expectedItems, items)
			mockRepo.AssertExpectations(t)
		})

		t.Run("Error on empty type", func(t *testing.T) {
			// Act
			items, err := catalogSvc.GetItems(ctx, "")

			// Assert
			assert.Error(t, err)
			assert.Nil(t, items)
			assert.Contains(t, err.Error(), "тип элемента каталога не может быть пустым")
		})
	})
}
