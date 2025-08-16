package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/farm/models"
)

// --- Greenhouse Methods ---

func (s *service) CreateGreenhouse(ctx context.Context, name, typeName string, landParcelID uuid.UUID) (*models.Greenhouse, error) {
	// Проверяем, существует ли земельный участок
	_, err := s.repo.GetLandParcelByID(ctx, landParcelID)
	if err != nil {
		return nil, fmt.Errorf("земельный участок с ID %s не найден: %w", landParcelID, err)
	}

	now := time.Now()
	greenhouse := &models.Greenhouse{
		ID:           uuid.New(),
		Name:         name,
		Type:         typeName,
		LandParcelID: landParcelID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = s.repo.CreateGreenhouse(ctx, greenhouse)
	if err != nil {
		return nil, err
	}

	return greenhouse, nil
}

func (s *service) GetGreenhouseByID(ctx context.Context, id uuid.UUID) (*models.Greenhouse, error) {
	return s.repo.GetGreenhouseByID(ctx, id)
}

func (s *service) GetGreenhousesByLandParcel(ctx context.Context, landParcelID uuid.UUID) ([]models.Greenhouse, error) {
	return s.repo.GetGreenhousesByLandParcel(ctx, landParcelID)
}

func (s *service) UpdateGreenhouse(ctx context.Context, id uuid.UUID, name, typeName string) (*models.Greenhouse, error) {
	greenhouse, err := s.repo.GetGreenhouseByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("теплица для обновления не найдена: %w", err)
	}

	greenhouse.Name = name
	greenhouse.Type = typeName
	greenhouse.UpdatedAt = time.Now()

	err = s.repo.UpdateGreenhouse(ctx, greenhouse)
	if err != nil {
		return nil, err
	}

	return greenhouse, nil
}

func (s *service) DeleteGreenhouse(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteGreenhouse(ctx, id)
}
