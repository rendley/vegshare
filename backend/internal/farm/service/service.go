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

	// LandParcel methods
	CreateLandParcel(ctx context.Context, name string, regionID uuid.UUID) (*models.LandParcel, error)
	GetLandParcelByID(ctx context.Context, id uuid.UUID) (*models.LandParcel, error)
	GetLandParcelsByRegion(ctx context.Context, regionID uuid.UUID) ([]models.LandParcel, error)
	UpdateLandParcel(ctx context.Context, id uuid.UUID, name string) (*models.LandParcel, error)
	DeleteLandParcel(ctx context.Context, id uuid.UUID) error
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
	region, err := s.repo.GetRegionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("регион для обновления не найден: %w", err)
	}

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

// --- LandParcel Methods ---

func (s *service) CreateLandParcel(ctx context.Context, name string, regionID uuid.UUID) (*models.LandParcel, error) {
	// Проверяем, существует ли регион
	_, err := s.repo.GetRegionByID(ctx, regionID)
	if err != nil {
		return nil, fmt.Errorf("регион с ID %s не найден: %w", regionID, err)
	}

	now := time.Now()
	parcel := &models.LandParcel{
		ID:        uuid.New(),
		Name:      name,
		RegionID:  regionID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.repo.CreateLandParcel(ctx, parcel)
	if err != nil {
		return nil, err
	}

	return parcel, nil
}

func (s *service) GetLandParcelByID(ctx context.Context, id uuid.UUID) (*models.LandParcel, error) {
	return s.repo.GetLandParcelByID(ctx, id)
}

func (s *service) GetLandParcelsByRegion(ctx context.Context, regionID uuid.UUID) ([]models.LandParcel, error) {
	return s.repo.GetLandParcelsByRegion(ctx, regionID)
}

func (s *service) UpdateLandParcel(ctx context.Context, id uuid.UUID, name string) (*models.LandParcel, error) {
	parcel, err := s.repo.GetLandParcelByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("земельный участок для обновления не найден: %w", err)
	}

	parcel.Name = name
	parcel.UpdatedAt = time.Now()

	err = s.repo.UpdateLandParcel(ctx, parcel)
	if err != nil {
		return nil, err
	}

	return parcel, nil
}

func (s *service) DeleteLandParcel(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteLandParcel(ctx, id)
}