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

// AuthRepository отвечает за взаимодействие с данными аутентификации в PostgreSQL.
type AuthRepository struct {
	db *sqlx.DB // Подключение к базе данных
}

// NewAuthRepository создаёт репозиторий экземляр с подключением к БД.
func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateUser создаёт нового пользователя в базе данных.
//   - ctx context.Context - контекст для контроля времени выполнения
//   - user *models.User - указатель на структуру пользователя
//   - error - ошибку, если операция не удалась
func (r *AuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	// Проверяем, что переданный пользователь не nil
	if user == nil {
		return errors.New("User connot be nil")
	}

	// Генерируем ID если не установлен
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	// SQL-запрос для вставки нового пользователя
	query := `
		INSERT INTO users (id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	// Устанавливаем временные метки
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Выполняем запрос с контекстом
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// GetUserByEmail находит пользователя по email.
// Принимает:
//   - ctx context.Context - контекст для контроля времени выполнения
//   - email string - email пользователя для поиска
//
// Возвращает:
//   - *models.User - найденного пользователя
//   - error - ошибку, если операция не удалась
func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	// Проверяем что email не пустой
	if email == "" {
		return nil, errors.New("Email cannot be empty")
	}

	// SQL-запрос для поиска пользователя
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`
	// Подготавливаем структуру для результата
	var user models.User

	// Выполняем запрос и сканируем результат
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return &user, nil
}

// UserExists проверяет существование пользователя с указанным email.
// Принимает:
//   - ctx context.Context - контекст для контроля времени выполнения
//   - email string - email для проверки
//
// Возвращает:
//   - bool - true если пользователь существует
//   - error - ошибку, если операция не удалась
func (r *AuthRepository) UserExists(ctx context.Context, email string) (bool, error) {
	// Проверяем, что email не пустой
	if email == "" {
		return false, errors.New("email cannot be empty")
	}

	// SQL-запрос для проверки существования
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`

	var exists bool

	// Выполняем запрос
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// SaveRefreshToken сохраняет refresh-токен в базе данных.
// Принимает:
//   - ctx context.Context - контекст для контроля времени выполнения
//   - userID uuid.UUID - ID пользователя
//   - token string - refresh-токен
//
// Возвращает:
//   - error - ошибку, если операция не удалась
func (r *AuthRepository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, token string) error {
	// Проверяем валидность входных данных
	if userID == uuid.Nil {
		return errors.New("invalid user ID")
	}
	if token == "" {
		return errors.New("token cannot be empty")
	}

	// SQL-запрос для вставки токена
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	// Генерируем UUID для токена
	tokenID := uuid.New()
	// Устанавливаем срок действия (30 дней)
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	// Текущее время
	createdAt := time.Now()

	// Выполняем запрос
	_, err := r.db.ExecContext(ctx, query,
		tokenID,
		userID,
		token,
		expiresAt,
		createdAt,
	)

	return err
}
