package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/farm/models"
)

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

func (s *service) RestoreLandParcel(ctx context.Context, id uuid.UUID) error {
	return s.repo.RestoreLandParcel(ctx, id)
}

func (s *service) GetAllLandParcelsIncludingDeleted(ctx context.Context) ([]models.LandParcel, error) {
	return s.repo.GetAllLandParcelsIncludingDeleted(ctx)
}
