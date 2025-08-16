// Пакет repository отвечает за прямой доступ к базе данных для сущностей фермы.
package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rendley/vegshare/backend/internal/farm/models"
)

// Repository - это ИНТЕРФЕЙС, который определяет "контракт" для нашего хранилища.
type Repository interface {
	GetAllCrops(ctx context.Context) ([]models.Crop, error)
}

// repository - это СТРУКТУРА, которая реализует интерфейс Repository.
type repository struct {
	db *sqlx.DB // Подключение к базе данных.
}

// NewRepository - это функция-конструктор.
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

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