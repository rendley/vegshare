package service

import (
	"context"
	"fmt"

	"github.com/rendley/vegshare/backend/internal/farm/models"
)

// --- Crop Methods ---

func (s *service) GetAllCrops(ctx context.Context) ([]models.Crop, error) {
	crops, err := s.repo.GetAllCrops(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить культуры в сервисе: %w", err)
	}
	return crops, nil
}
