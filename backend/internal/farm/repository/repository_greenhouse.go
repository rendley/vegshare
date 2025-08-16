package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/farm/models"
)

// --- Greenhouse Methods ---

func (r *repository) CreateGreenhouse(ctx context.Context, greenhouse *models.Greenhouse) error {
	query := `INSERT INTO greenhouses (id, land_parcel_id, name, type, created_at, updated_at) VALUES (:id, :land_parcel_id, :name, :type, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, greenhouse)
	if err != nil {
		return fmt.Errorf("не удалось создать теплицу: %w", err)
	}
	return nil
}

func (r *repository) GetGreenhouseByID(ctx context.Context, id uuid.UUID) (*models.Greenhouse, error) {
	var greenhouse models.Greenhouse
	query := `SELECT * FROM greenhouses WHERE id = $1`
	err := r.db.GetContext(ctx, &greenhouse, query, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить теплицу по ID: %w", err)
	}
	return &greenhouse, nil
}

func (r *repository) GetGreenhousesByLandParcel(ctx context.Context, landParcelID uuid.UUID) ([]models.Greenhouse, error) {
	var greenhouses []models.Greenhouse
	query := `SELECT * FROM greenhouses WHERE land_parcel_id = $1`
	err := r.db.SelectContext(ctx, &greenhouses, query, landParcelID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список теплиц для земельного участка: %w", err)
	}
	return greenhouses, nil
}

func (r *repository) UpdateGreenhouse(ctx context.Context, greenhouse *models.Greenhouse) error {
	query := `UPDATE greenhouses SET name = :name, type = :type, updated_at = :updated_at WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, greenhouse)
	if err != nil {
		return fmt.Errorf("не удалось обновить теплицу: %w", err)
	}
	return nil
}

func (r *repository) DeleteGreenhouse(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM greenhouses WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить теплицу: %w", err)
	}
	return nil
}
