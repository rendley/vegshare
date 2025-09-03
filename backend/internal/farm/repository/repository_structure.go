package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/farm/models"
)

// --- Structure Methods ---

func (r *repository) CreateStructure(ctx context.Context, structure *models.Structure) error {
	query := `INSERT INTO structures (id, land_parcel_id, name, type, created_at, updated_at) VALUES (:id, :land_parcel_id, :name, :type, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, structure)
	if err != nil {
		return fmt.Errorf("не удалось создать строение: %w", err)
	}
	return nil
}

func (r *repository) GetStructureByID(ctx context.Context, id uuid.UUID) (*models.Structure, error) {
	var structure models.Structure
	query := `SELECT * FROM structures WHERE id = $1`
	err := r.db.GetContext(ctx, &structure, query, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить строение по ID: %w", err)
	}
	return &structure, nil
}

func (r *repository) GetStructuresByLandParcel(ctx context.Context, landParcelID uuid.UUID) ([]models.Structure, error) {
	var structures []models.Structure
	query := `SELECT * FROM structures WHERE land_parcel_id = $1`
	err := r.db.SelectContext(ctx, &structures, query, landParcelID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список строений для земельного участка: %w", err)
	}
	return structures, nil
}

func (r *repository) UpdateStructure(ctx context.Context, structure *models.Structure) error {
	query := `UPDATE structures SET name = :name, type = :type, updated_at = :updated_at WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, structure)
	if err != nil {
		return fmt.Errorf("не удалось обновить строение: %w", err)
	}
	return nil
}

func (r *repository) DeleteStructure(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM structures WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить строение: %w", err)
	}
	return nil
}

func (r *repository) GetStructureTypes(ctx context.Context) ([]string, error) {
	var types []string
	query := `SELECT DISTINCT type FROM structures ORDER BY type`
	err := r.db.SelectContext(ctx, &types, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить типы строений: %w", err)
	}
	return types, nil
}