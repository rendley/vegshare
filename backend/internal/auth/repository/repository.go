package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rendley/vegshare/backend/internal/auth/models"
	"time"
)

// AuthRepository defines the interface for authentication data access.
type AuthRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UserExists(ctx context.Context, email string) (bool, error)
	SaveRefreshToken(ctx context.Context, userID uuid.UUID, token string) error
}

// authRepository is the implementation of AuthRepository.
type authRepository struct {
	db *sqlx.DB
}

// NewAuthRepository creates a new instance of AuthRepository.
func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepository{db: db}
}

// CreateUser creates a new user in the database.
func (r *authRepository) CreateUser(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	query := `
		INSERT INTO users (id, name, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// GetUserByEmail finds a user by email.
func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	query := `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`
	var user models.User

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrInvalidCredentials
		}
		return nil, err
	}
	return &user, nil
}

// UserExists checks if a user with the given email exists.
func (r *authRepository) UserExists(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, errors.New("email cannot be empty")
	}

	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`

	var exists bool

	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// SaveRefreshToken saves a refresh token to the database.
func (r *authRepository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, token string) error {
	if userID == uuid.Nil {
		return errors.New("invalid user ID")
	}
	if token == "" {
		return errors.New("token cannot be empty")
	}

	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	tokenID := uuid.New()
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	createdAt := time.Now()

	_, err := r.db.ExecContext(ctx, query,
		tokenID,
		userID,
		token,
		expiresAt,
		createdAt,
	)

	return err
}
