package database

import (
	"fmt"
	"github.com/jmoiron/sqlx" // <-- ИЗМЕНЕНО
	_ "github.com/lib/pq"
	"github.com/rendley/vegshare/backend/pkg/config"
)

// New создает подключение к PostgreSQL и возвращает *sqlx.DB
func New(cfg *config.Config) (*sqlx.DB, error) { // <-- ИЗМЕНЕНО
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// Используем sqlx.Open для подключения
	db, err := sqlx.Open("postgres", connStr) // <-- ИЗМЕНЕНО
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
