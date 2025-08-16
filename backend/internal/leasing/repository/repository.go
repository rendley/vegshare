// Пакет repository отвечает за доступ к данным, связанным с арендой.
package repository

import "github.com/jmoiron/sqlx"

// Repository определяет контракт для хранилища данных аренды.
type Repository interface {
	// Здесь будут методы для работы с таблицей plot_leases
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}
