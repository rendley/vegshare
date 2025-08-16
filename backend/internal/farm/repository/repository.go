// Пакет repository отвечает за прямой доступ к базе данных для сущностей фермы.
package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rendley/vegshare/backend/internal/farm/models"
)

// Repository - это ИНТЕРФЕЙС, который определяет "контракт" для нашего хранилища.
type Repository interface {
	// Crop methods
	GetAllCrops(ctx context.Context) ([]models.Crop, error)

	// Region methods
	CreateRegion(ctx context.Context, region *models.Region) error
	GetRegionByID(ctx context.Context, id uuid.UUID) (*models.Region, error)
	GetAllRegions(ctx context.Context) ([]models.Region, error)
	UpdateRegion(ctx context.Context, region *models.Region) error
	DeleteRegion(ctx context.Context, id uuid.UUID) error
}

// repository - это СТРУКТУРА, которая реализует интерфейс Repository.
type repository struct {
	db *sqlx.DB // Подключение к базе данных.
}

// NewRepository - это функция-конструктор.
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

// --- Crop Methods ---

// GetAllCrops - это реализация метода для получения всех культур из базы данных.
func (r *repository) GetAllCrops(ctx context.Context) ([]models.Crop, error) {
	var crops []models.Crop
	query := "SELECT * FROM crops"
	err := r.db.SelectContext(ctx, &crops, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список культур: %w", err)
	}
	return crops, nil
}

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
	query := `SELECT * FROM regions WHERE id = $1`
	err := r.db.GetContext(ctx, &region, query, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить регион по ID: %w", err)
	}
	return &region, nil
}

func (r *repository) GetAllRegions(ctx context.Context) ([]models.Region, error) {
	var regions []models.Region
	query := `SELECT * FROM regions`
	err := r.db.SelectContext(ctx, &regions, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список регионов: %w", err)
	}
	return regions, nil
}

func (r *repository) UpdateRegion(ctx context.Context, region *models.Region) error {
	query := `UPDATE regions SET name = :name, updated_at = :updated_at WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, region)
	if err != nil {
		return fmt.Errorf("не удалось обновить регион: %w", err)
	}
	return nil
}

func (r *repository) DeleteRegion(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM regions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить регион: %w", err)
	}
	return nil
}
