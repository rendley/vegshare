package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/farm/models"
)

// --- Structure Methods ---

func (s *service) CreateStructure(ctx context.Context, name, typeName string, landParcelID uuid.UUID) (*models.Structure, error) {
	// Проверяем, существует ли земельный участок
	_, err := s.repo.GetLandParcelByID(ctx, landParcelID)
	if err != nil {
		return nil, fmt.Errorf("земельный участок с ID %s не найден: %w", landParcelID, err)
	}

	now := time.Now()
	structure := &models.Structure{
		ID:           uuid.New(),
		Name:         name,
		Type:         typeName,
		LandParcelID: landParcelID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = s.repo.CreateStructure(ctx, structure)
	if err != nil {
		return nil, err
	}

	return structure, nil
}

func (s *service) GetStructureByID(ctx context.Context, id uuid.UUID) (*models.Structure, error) {
	return s.repo.GetStructureByID(ctx, id)
}

func (s *service) GetStructuresByLandParcel(ctx context.Context, landParcelID uuid.UUID) ([]models.Structure, error) {
	return s.repo.GetStructuresByLandParcel(ctx, landParcelID)
}

func (s *service) UpdateStructure(ctx context.Context, id uuid.UUID, name, typeName string) (*models.Structure, error) {
	structure, err := s.repo.GetStructureByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("строени для обновления не найдено: %w", err)
	}

	structure.Name = name
	structure.Type = typeName
	structure.UpdatedAt = time.Now()

	err = s.repo.UpdateStructure(ctx, structure)
	if err != nil {
		return nil, err
	}

	return structure, nil
}

func (s *service) DeleteStructure(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteStructure(ctx, id)
}

func (s *service) GetStructureTypes(ctx context.Context) ([]string, error) {
	return s.repo.GetStructureTypes(ctx)
}