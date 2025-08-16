// Пакет service содержит бизнес-логику, связанную с фермами.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/farm/models"
	"github.com/rendley/vegshare/backend/internal/farm/repository"
)

// Service - это интерфейс, определяющий контракт для сервиса фермы.
type Service interface {
	// Crop methods
	GetAllCrops(ctx context.Context) ([]models.Crop, error)

	// Region methods
	CreateRegion(ctx context.Context, name string) (*models.Region, error)
	GetRegionByID(ctx context.Context, id uuid.UUID) (*models.Region, error)
	GetAllRegions(ctx context.Context) ([]models.Region, error)
	UpdateRegion(ctx context.Context, id uuid.UUID, name string) (*models.Region, error)
	DeleteRegion(ctx context.Context, id uuid.UUID) error
}

// service - это приватная структура, реализующая интерфейс Service.
type service struct {
	repo repository.Repository
}

// NewFarmService - это конструктор для нашего сервиса.
func NewFarmService(repo repository.Repository) Service {
	return &service{repo: repo}
}

// --- Crop Methods ---

func (s *service) GetAllCrops(ctx context.Context) ([]models.Crop, error) {
	crops, err := s.repo.GetAllCrops(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить культуры в сервисе: %w", err)
	}
	return crops, nil
}

// --- Region Methods ---

func (s *service) CreateRegion(ctx context.Context, name string) (*models.Region, error) {
	now := time.Now()
	region := &models.Region{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := s.repo.CreateRegion(ctx, region)
	if err != nil {
		return nil, err
	}

	return region, nil
}

func (s *service) GetRegionByID(ctx context.Context, id uuid.UUID) (*models.Region, error) {
	return s.repo.GetRegionByID(ctx, id)
}

func (s *service) GetAllRegions(ctx context.Context) ([]models.Region, error) {
	return s.repo.GetAllRegions(ctx)
}

func (s *service) UpdateRegion(ctx context.Context, id uuid.UUID, name string) (*models.Region, error) {
	// Сначала получаем регион, чтобы убедиться, что он существует
	region, err := s.repo.GetRegionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("регион для обновления не найден: %w", err)
	}

	// Обновляем поля
	region.Name = name
	region.UpdatedAt = time.Now()

	err = s.repo.UpdateRegion(ctx, region)
	if err != nil {
		return nil, err
	}

	return region, nil
}

func (s *service) DeleteRegion(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteRegion(ctx, id)
}
