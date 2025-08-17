package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rendley/vegshare/backend/internal/farm/models"
)

// Repository определяет контракт для хранилища.
type Repository interface {
	CreatePlotCrop(ctx context.Context, plotCrop *models.PlotCrop) error
	GetPlotCrops(ctx context.Context, plotID uuid.UUID) ([]models.PlotCrop, error)
	DeletePlotCrop(ctx context.Context, plantingID uuid.UUID) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreatePlotCrop(ctx context.Context, plotCrop *models.PlotCrop) error {
	query := `INSERT INTO plot_crops (id, plot_id, crop_id, lease_id, planted_at, status, created_at, updated_at) VALUES (:id, :plot_id, :crop_id, :lease_id, :planted_at, :status, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, plotCrop)
	if err != nil {
		return fmt.Errorf("не удалось создать запись о посадке: %w", err)
	}
	return nil
}

func (r *repository) GetPlotCrops(ctx context.Context, plotID uuid.UUID) ([]models.PlotCrop, error) {
	var plotCrops []models.PlotCrop
	query := `SELECT * FROM plot_crops WHERE plot_id = $1`
	err := r.db.SelectContext(ctx, &plotCrops, query, plotID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить посадки для грядки: %w", err)
	}
	return plotCrops, nil
}

func (r *repository) DeletePlotCrop(ctx context.Context, plantingID uuid.UUID) error {
	query := `DELETE FROM plot_crops WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, plantingID)
	if err != nil {
		return fmt.Errorf("не удалось удалить запись о посадке: %w", err)
	}
	return nil
}
