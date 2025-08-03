package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Драйвер PostgreSQL
	"github.com/rendley/backend/pkg/config"
)

// New создает подключение к PostgreSQL
func New(cfg *config.Config) (*sql.DB, error) {
	// Формируем строку подключения
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)
	// Открытие соединения
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}
