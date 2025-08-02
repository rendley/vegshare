package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// AuthRepository работает с данными пользователей в PostgreSQL.
type AuthRepository struct {
	db *sql.DB
}

// NewAuthRepository создаёт репозиторий с подключением к БД.
func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(email, hashedPassword string) (string, error) {
	// 1. Проверяем валидность email
	if email == "" {
		return "", errors.New("Email connot be empty")
	}
	// 2. SQL-запрос с RETURNING (PostgreSQL)
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`
	// 3. Контекст с таймаутом (5 секунд)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 4. Выполняем запрос
	var userID string
	err := r.db.QueryRowContext(ctx, query, email, hashedPassword).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}
