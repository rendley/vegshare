package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/farm/models"
)

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
