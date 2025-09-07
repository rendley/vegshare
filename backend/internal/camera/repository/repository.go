package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rendley/vegshare/backend/internal/camera/models"
)

// Repository defines the contract for camera data storage.
type Repository interface {
	CreateCamera(ctx context.Context, camera *models.Camera) error
	GetCamerasByUnitID(ctx context.Context, unitID uuid.UUID, unitType string) ([]models.Camera, error)
	GetCameraByID(ctx context.Context, cameraID uuid.UUID) (*models.Camera, error)
	DeleteCamera(ctx context.Context, cameraID uuid.UUID) error
}

// repository implements the Repository interface.
type repository struct {
	db *sqlx.DB
}

// NewRepository is a constructor for the repository.
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateCamera(ctx context.Context, camera *models.Camera) error {
	query := `INSERT INTO cameras (id, unit_id, unit_type, name, rtsp_path_name, created_at, updated_at)
		VALUES (:id, :unit_id, :unit_type, :name, :rtsp_path_name, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, camera)
	if err != nil {
		return fmt.Errorf("не удалось создать камеру: %w", err)
	}
	return nil
}

func (r *repository) GetCamerasByUnitID(ctx context.Context, unitID uuid.UUID, unitType string) ([]models.Camera, error) {
	var cameras []models.Camera
	query := `SELECT * FROM cameras WHERE unit_id = $1 AND unit_type = $2`
	err := r.db.SelectContext(ctx, &cameras, query, unitID, unitType)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список камер для юнита: %w", err)
	}
	return cameras, nil
}

func (r *repository) DeleteCamera(ctx context.Context, cameraID uuid.UUID) error {
	query := `DELETE FROM cameras WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, cameraID)
	if err != nil {
		return fmt.Errorf("не удалось удалить камеру: %w", err)
	}
	return nil
}

func (r *repository) GetCameraByID(ctx context.Context, cameraID uuid.UUID) (*models.Camera, error) {
	var camera models.Camera
	query := `SELECT * FROM cameras WHERE id = $1`
	err := r.db.GetContext(ctx, &camera, query, cameraID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить камеру по ID: %w", err)
	}
	return &camera, nil
}
