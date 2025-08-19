package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rendley/vegshare/backend/internal/catalog/models"
)

// Repository defines the contract for the catalog data storage.
type Repository interface {
	GetAllCrops(ctx context.Context) ([]models.Crop, error)
	GetCropByID(ctx context.Context, id uuid.UUID) (*models.Crop, error)
	CreateCrop(ctx context.Context, crop *models.Crop) error
}

// repository implements the Repository interface.
type repository struct {
	db *sqlx.DB
}

// NewRepository is a constructor for the repository.
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAllCrops(ctx context.Context) ([]models.Crop, error) {
	var crops []models.Crop
	query := "SELECT * FROM crops"
	err := r.db.SelectContext(ctx, &crops, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список культур: %w", err)
	}
	return crops, nil
}

func (r *repository) GetCropByID(ctx context.Context, id uuid.UUID) (*models.Crop, error) {
	var crop models.Crop
	query := `SELECT * FROM crops WHERE id = $1`
	err := r.db.GetContext(ctx, &crop, query, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить культуру по ID: %w", err)
	}
	return &crop, nil
}

func (r *repository) CreateCrop(ctx context.Context, crop *models.Crop) error {
	query := `INSERT INTO crops (id, name, description, planting_time, harvest_time, created_at, updated_at) VALUES (:id, :name, :description, :planting_time, :harvest_time, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, crop)
	if err != nil {
		return fmt.Errorf("не удалось создать культуру: %w", err)
	}
	return nil
}
