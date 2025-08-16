package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/farm/models"
)

// --- LandParcel Methods ---

func (r *repository) CreateLandParcel(ctx context.Context, parcel *models.LandParcel) error {
	query := `INSERT INTO land_parcels (id, region_id, name, created_at, updated_at) VALUES (:id, :region_id, :name, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, parcel)
	if err != nil {
		return fmt.Errorf("не удалось создать земельный участок: %w", err)
	}
	return nil
}

func (r *repository) GetLandParcelByID(ctx context.Context, id uuid.UUID) (*models.LandParcel, error) {
	var parcel models.LandParcel
	query := `SELECT * FROM land_parcels WHERE id = $1`
	err := r.db.GetContext(ctx, &parcel, query, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить земельный участок по ID: %w", err)
	}
	return &parcel, nil
}

func (r *repository) GetLandParcelsByRegion(ctx context.Context, regionID uuid.UUID) ([]models.LandParcel, error) {
	var parcels []models.LandParcel
	query := `SELECT * FROM land_parcels WHERE region_id = $1`
	err := r.db.SelectContext(ctx, &parcels, query, regionID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список земельных участков для региона: %w", err)
	}
	return parcels, nil
}

func (r *repository) UpdateLandParcel(ctx context.Context, parcel *models.LandParcel) error {
	query := `UPDATE land_parcels SET name = :name, updated_at = :updated_at WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, parcel)
	if err != nil {
		return fmt.Errorf("не удалось обновить земельный участок: %w", err)
	}
	return nil
}

func (r *repository) DeleteLandParcel(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM land_parcels WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить земельный участок: %w", err)
	}
	return nil
}
