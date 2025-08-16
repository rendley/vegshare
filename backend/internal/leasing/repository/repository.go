// Пакет repository отвечает за доступ к данным, связанным с арендой.
package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rendley/vegshare/backend/internal/farm/models"
)

// Repository определяет контракт для хранилища данных аренды.
type Repository interface {
	CreateLease(ctx context.Context, lease *models.PlotLease) error
	GetLeasesByUserID(ctx context.Context, userID uuid.UUID) ([]models.PlotLease, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateLease(ctx context.Context, lease *models.PlotLease) error {
	query := `INSERT INTO plot_leases (id, plot_id, user_id, start_date, end_date, status, created_at, updated_at) VALUES (:id, :plot_id, :user_id, :start_date, :end_date, :status, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, lease)
	if err != nil {
		return fmt.Errorf("не удалось создать запись аренды: %w", err)
	}
	return nil
}

func (r *repository) GetLeasesByUserID(ctx context.Context, userID uuid.UUID) ([]models.PlotLease, error) {
	var leases []models.PlotLease
	query := `SELECT * FROM plot_leases WHERE user_id = $1 AND status = 'active'`
	err := r.db.SelectContext(ctx, &leases, query, userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список аренд для пользователя: %w", err)
	}
	return leases, nil
}
