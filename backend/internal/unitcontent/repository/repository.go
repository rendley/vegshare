package repository

import (
	"context"
	"github.com/rendley/vegshare/backend/internal/unitcontent/models"
	"github.com/rendley/vegshare/backend/pkg/database"
)

// Repository определяет интерфейс для взаимодействия с хранилищем содержимого юнитов.
type Repository interface {
	// CreateOrUpdate создает новую запись или обновляет существующую, если контент для юнита уже есть.
	CreateOrUpdate(ctx context.Context, content *models.UnitContent) error
}

// repository - реализация Repository для PostgreSQL.
type repository struct {
	db database.DBTX
}

// NewRepository - конструктор для repository.
func NewRepository(db database.DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateOrUpdate(ctx context.Context, content *models.UnitContent) error {
	// Используем мощный синтаксис PostgreSQL "ON CONFLICT ... DO UPDATE"
	// Он позволяет атомарно выполнить INSERT или UPDATE.
	query := `INSERT INTO unit_contents (id, unit_id, unit_type, item_id, quantity, updated_at)
			 VALUES (:id, :unit_id, :unit_type, :item_id, :quantity, :updated_at)
			 ON CONFLICT (unit_id, unit_type) DO UPDATE
			 SET item_id = EXCLUDED.item_id,
				 quantity = EXCLUDED.quantity,
				 updated_at = EXCLUDED.updated_at`

	_, err := r.db.NamedExecContext(ctx, query, content)
	return err
}