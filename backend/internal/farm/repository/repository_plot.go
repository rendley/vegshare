package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/farm/models"
)

// --- Plot Methods ---

func (r *repository) CreatePlot(ctx context.Context, plot *models.Plot) error {
	query := `INSERT INTO plots (id, greenhouse_id, name, size, status, created_at, updated_at) VALUES (:id, :greenhouse_id, :name, :size, :status, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, plot)
	if err != nil {
		return fmt.Errorf("не удалось создать грядку: %w", err)
	}
	return nil
}

func (r *repository) GetPlotByID(ctx context.Context, id uuid.UUID) (*models.Plot, error) {
	var plot models.Plot
	query := `SELECT * FROM plots WHERE id = $1`
	err := r.db.GetContext(ctx, &plot, query, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить грядку по ID: %w", err)
	}
	return &plot, nil
}

func (r *repository) GetPlotsByGreenhouse(ctx context.Context, greenhouseID uuid.UUID) ([]models.Plot, error) {
	var plots []models.Plot
	query := `SELECT * FROM plots WHERE greenhouse_id = $1`
	err := r.db.SelectContext(ctx, &plots, query, greenhouseID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список грядок для теплицы: %w", err)
	}
	return plots, nil
}

func (r *repository) UpdatePlot(ctx context.Context, plot *models.Plot) error {
	query := `UPDATE plots SET name = :name, size = :size, status = :status, updated_at = :updated_at WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, plot)
	if err != nil {
		return fmt.Errorf("не удалось обновить грядку: %w", err)
	}
	return nil
}

func (r *repository) DeletePlot(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM plots WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить грядку: %w", err)
	}
	return nil
}
