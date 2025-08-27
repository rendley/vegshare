package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rendley/vegshare/backend/internal/user/models"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserProfile, error)
	GetAllUsers(ctx context.Context) ([]models.UserProfile, error)
	UpdateUser(ctx context.Context, user *models.UserProfile) error
	UpdateUserRole(ctx context.Context, id uuid.UUID, role string) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// userRepository is the implementation of UserRepository.
type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

// GetUserByID retrieves a user by their ID.
func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserProfile, error) {
	var user models.UserProfile
	query := "SELECT id, email, name, role, avatar_url, created_at, updated_at FROM users WHERE id = $1"
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось найти пользователя по id: %w", err)
	}
	return &user, nil
}

// GetAllUsers retrieves all users from the database.
func (r *userRepository) GetAllUsers(ctx context.Context) ([]models.UserProfile, error) {
	var users []models.UserProfile
	query := "SELECT id, email, name, role, avatar_url, created_at, updated_at FROM users ORDER BY created_at DESC"
	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список пользователей: %w", err)
	}
	return users, nil
}

// UpdateUser updates a user's profile.
func (r *userRepository) UpdateUser(ctx context.Context, user *models.UserProfile) error {
	query := `UPDATE users SET 
				name = :name, 
				avatar_url = :avatar_url,
				updated_at = now()
			  WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("не удалось обновить пользователя: %w", err)
	}
	return nil
}

// UpdateUserRole updates a user's role.
func (r *userRepository) UpdateUserRole(ctx context.Context, id uuid.UUID, role string) error {
	query := "UPDATE users SET role = $1, updated_at = now() WHERE id = $2"
	result, err := r.db.ExecContext(ctx, query, role, id)
	if err != nil {
		return fmt.Errorf("не удалось обновить роль пользователя: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось проверить результат обновления: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("пользователь с ID %s не найден для обновления роли", id)
	}
	return nil
}

// DeleteUser deletes a user by their ID.
func (r *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить пользователя: %w", err)
	}
	return nil
}