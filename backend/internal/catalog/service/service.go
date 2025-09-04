package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/catalog/models"
	"github.com/rendley/vegshare/backend/internal/catalog/repository"
)

// Service defines the contract for the catalog service.
type Service interface {
	GetItems(ctx context.Context, itemType string) ([]models.CatalogItem, error)
	GetItemByID(ctx context.Context, id uuid.UUID) (*models.CatalogItem, error)
	CreateItem(ctx context.Context, itemType, name, description string, attributes models.JSONB) (*models.CatalogItem, error)
}

// service implements the Service interface.
type service struct {
	repo repository.Repository
}

// NewService is a constructor for the service.
func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetItems(ctx context.Context, itemType string) ([]models.CatalogItem, error) {
	if itemType == "" {
		return nil, fmt.Errorf("тип элемента каталога не может быть пустым")
	}
	items, err := s.repo.GetItemsByType(ctx, itemType)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить элементы каталога в сервисе: %w", err)
	}
	return items, nil
}

func (s *service) GetItemByID(ctx context.Context, id uuid.UUID) (*models.CatalogItem, error) {
	return s.repo.GetItemByID(ctx, id)
}

func (s *service) CreateItem(ctx context.Context, itemType, name, description string, attributes models.JSONB) (*models.CatalogItem, error) {
	now := time.Now()
	item := &models.CatalogItem{
		ID:          uuid.New(),
		ItemType:    itemType,
		Name:        name,
		Description: description,
		Attributes:  attributes,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.CreateItem(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}