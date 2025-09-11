package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/farm/models"
)

// --- Region Methods ---

func (r *repository) CreateRegion(ctx context.Context, region *models.Region) error {
	query := `INSERT INTO regions (id, name, created_at, updated_at) VALUES (:id, :name, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, region)
	if err != nil {
		return fmt.Errorf("не удалось создать регион: %w", err)
	}
	return nil
}

func (r *repository) GetRegionByID(ctx context.Context, id uuid.UUID) (*models.Region, error) {
	var region models.Region
	query := `SELECT * FROM regions WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &region, query, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить регион по ID: %w", err)
	}
	return &region, nil
}

func (r *repository) GetAllRegions(ctx context.Context) ([]models.Region, error) {
	var regions []models.Region
	query := `SELECT * FROM regions WHERE deleted_at IS NULL`
	err := r.db.SelectContext(ctx, &regions, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список регионов: %w", err)
	}
	return regions, nil
}

func (r *repository) UpdateRegion(ctx context.Context, region *models.Region) error {
	query := `UPDATE regions SET name = :name, updated_at = :updated_at WHERE id = :id AND deleted_at IS NULL`
	_, err := r.db.NamedExecContext(ctx, query, region)
	if err != nil {
		return fmt.Errorf("не удалось обновить регион: %w", err)
	}
	return nil
}

func (r *repository) DeleteRegion(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE regions SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить регион: %w", err)
	}
	return nil
}

func (r *repository) RestoreRegion(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE regions SET deleted_at = NULL WHERE id = $1 AND deleted_at IS NOT NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось восстановить регион: %w", err)
	}
	return nil
}

func (r *repository) GetAllRegionsIncludingDeleted(ctx context.Context) ([]models.Region, error) {
	var regions []models.Region
	query := `SELECT * FROM regions`
	err := r.db.SelectContext(ctx, &regions, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить полный список регионов: %w", err)
	}
	return regions, nil
}
