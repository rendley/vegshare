package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/plot/models"
	"github.com/rendley/vegshare/backend/pkg/database"
)

// Repository defines the contract for plot data storage.
type Repository interface {
	CreatePlot(ctx context.Context, plot *models.Plot) error
	GetPlotByID(ctx context.Context, id uuid.UUID) (*models.Plot, error)
	GetPlotsByStructure(ctx context.Context, structureID uuid.UUID) ([]models.Plot, error)
	UpdatePlot(ctx context.Context, plot *models.Plot) error
	DeletePlot(ctx context.Context, id uuid.UUID) error
}

// repository implements the Repository interface.
type repository struct {
	db database.DBTX
}

// NewRepository is a constructor for the repository.
func NewRepository(db database.DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreatePlot(ctx context.Context, plot *models.Plot) error {
	query := `INSERT INTO plots (id, structure_id, name, size, status, created_at, updated_at) VALUES (:id, :structure_id, :name, :size, :status, :created_at, :updated_at)`
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

func (r *repository) GetPlotsByStructure(ctx context.Context, structureID uuid.UUID) ([]models.Plot, error) {
	var plots []models.Plot
	query := `SELECT * FROM plots WHERE structure_id = $1`
	err := r.db.SelectContext(ctx, &plots, query, structureID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список грядок для строения: %w", err)
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
