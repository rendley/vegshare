// Пакет repository отвечает за доступ к данным, связанным с арендой.
package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/leasing/models"
	"github.com/rendley/vegshare/backend/pkg/database"
)

// Repository определяет контракт для хранилища данных аренды.
// Теперь он работает с универсальной моделью models.Lease.
type Repository interface {
	CreateLease(ctx context.Context, lease *models.Lease) error
	GetLeasesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Lease, error)
}

type repository struct {
	db database.DBTX
}

func NewRepository(db database.DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateLease(ctx context.Context, lease *models.Lease) error {
	// Запрос обновлен для работы с новой таблицей 'leases' и ее полями.
	query := `INSERT INTO leases (id, unit_id, unit_type, user_id, start_date, end_date, status, created_at, updated_at) 
	          VALUES (:id, :unit_id, :unit_type, :user_id, :start_date, :end_date, :status, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, lease)
	if err != nil {
		return fmt.Errorf("не удалось создать запись аренды: %w", err)
	}
	return nil
}

func (r *repository) GetLeasesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Lease, error) {
	// Используем новую универсальную модель Lease
	var leases []models.Lease
	// Запрос обновлен для работы с новой таблицей 'leases'.
	query := `SELECT * FROM leases WHERE user_id = $1 AND status = 'active'`
	err := r.db.SelectContext(ctx, &leases, query, userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список аренд для пользователя: %w", err)
	}
	return leases, nil
}
