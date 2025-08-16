package repository

import (
	"context"
	"fmt"

	"github.com/rendley/vegshare/backend/internal/farm/models"
)

// --- Crop Methods ---

func (r *repository) GetAllCrops(ctx context.Context) ([]models.Crop, error) {
	var crops []models.Crop
	query := "SELECT * FROM crops"
	err := r.db.SelectContext(ctx, &crops, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список культур: %w", err)
	}
	return crops, nil
}
