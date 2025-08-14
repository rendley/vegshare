// Пакет handler содержит базовую структуру и конструктор.
// В данном пакете также собраны методы обработчиков связанные с auth
package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rendley/backend/internal/auth/models"
	"github.com/rendley/backend/internal/auth/repository"
	"github.com/rendley/backend/pkg/api"
	"github.com/rendley/backend/pkg/jwt"
	"github.com/rendley/backend/pkg/security"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Handler — корневая структура для всех обработчиков.
type AuthHandler struct {
	db             *sql.DB                    // Подключение к PostgreSQL
	logger         *logrus.Logger             // Логгер
	passwordHasher security.PasswordHasher    // Хеширование пассводра
	repository     *repository.AuthRepository // репозиторий для работы с базой
	validate       *validator.Validate        // валидатор
	jwtGenerator   jwt.Generator              // генератор токенов
}

// New создаёт экземпляр Handler c зависимостями.
func NewAuthHandler(db *sql.DB, hasher security.PasswordHasher, logger *logrus.Logger, jwtGen jwt.Generator) *AuthHandler {
	return &AuthHandler{
		db:             db,
		logger:         logger,
		passwordHasher: hasher,
		repository:     repository.NewAuthRepository(db), // Инициализируем репозиторий
		validate:       validator.New(),                  // Инициализируем валидатор
		jwtGenerator:   jwtGen,                           // Инициализируем генератор
	}
}

//############################## Hendlers ###################################//

// homeHandler обрабатывает запросы к корневому пути ("/").
func (h *AuthHandler) homeHandler(w http.ResponseWriter, r *http.Request) {
	//Проверяем что запрос именно GET
	if r.Method != http.MethodGet {
		// Если метод не GET, возвращаем ошибку 405 (Method Not Allowed).
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
		return
	}
	// Пишем ответ в http.ResponseWriter.
	fmt.Fprintf(w, "Welcome to VegShare!")
}

// registerHandler обрабатывает POST запросы "/register"
func (h *AuthHandler) registerHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Логируем начало обработки запроса и добавляем контекст
	ctx := r.Context()
	h.logger.Info("Register request received")

	// 2. Проверяем метод запроса (должен быть POST)
	if r.Method != http.MethodPost {
		h.logger.Warn("Wrong method used for registration")
		api.RespondWithError(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	// 3. Читаем тело запроса делаем это из локальной структуры
	//  а если будут повторения всегда можно вынести в types.go
	var request RegisterRequest

	// 4. Парсим JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Warn("Failed to parse JSON: %v", err)
		api.RespondWithError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// 5. Валидация обязательных полей
	if err := h.validate.Struct(request); err != nil {
		h.logger.Warnf("Validation failed: %v", err)
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 6. Проверка существования пользователя
	exists, err := h.repository.UserExists(ctx, request.Email)
	if err != nil {
		h.logger.Errorf("Failed to ckeck user existence", http.StatusInternalServerError)
		api.RespondWithError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if exists {
		h.logger.Warnf("User already exists with email: %s", request.Email)
		api.RespondWithError(w, "User already exists", http.StatusConflict)
		return
	}

	// 7. Хешируем пароль перед сохранением в БД
	hashedPassword, err := h.passwordHasher.Hash(request.Password)
	if err != nil {
		h.logger.Errorf("Password hashing failed: %v", err)
		api.RespondWithError(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	// 8. Сохраняем пользователя в БД
	user := models.User{
		ID:           uuid.New(), // Явно устанавливаем ID
		Email:        request.Email,
		PasswordHash: hashedPassword,
	}

	if user.ID == uuid.Nil {
		h.logger.Error("Generated invalid UUID")
		api.RespondWithError(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	if err := h.repository.CreateUser(ctx, &user); err != nil {
		h.logger.Errorf("Failed to create user: %v", err)
		api.RespondWithError(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	// 8. Генерация токена
	tokens, err := h.generateTokens(user.ID)
	if err != nil {
		h.logger.Errorf("Token generrating failed: %v", err)
		api.RespondWithError(w, "Login failed", http.StatusInternalServerError)
		return
	}

	// Сохранение refresh-токена
	if err := h.repository.SaveRefreshToken(ctx, user.ID, tokens.RefreshToken); err != nil {
		h.logger.Errorf("Failed to save refresh token: %v", err)
		api.RespondWithError(w, "Login failed", http.StatusInternalServerError)
		return
	}

	// 9. Отправляем успешный ответ
	response := LoginResponse{
		UserID: user.ID.String(),
		TokenPair: TokenPair{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}
	api.RespondWithJSON(h.logger, w, response, http.StatusCreated)
}

// loginHandler обрабатывает POST /login.
func (h *AuthHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info("Login request received")

	if r.Method != http.MethodPost {
		h.respondWithError(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warnf("Failed to parse JSON: %v", err)
		h.respondWithError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Валидация
	if err := h.validate.Struct(req); err != nil {
		h.respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Поиск пользователя в БД
	user, err := h.repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		h.logger.Warnf("User not found: %v", err)
		h.respondWithError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Проверка пароля
	if !h.passwordHasher.Check(user.PasswordHash, req.Password) {
		h.logger.Warn("Invalid password attempt")
		h.respondWithError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Генерация токенов
	tokens, err := h.generateTokens(user.ID)
	if err != nil {
		h.logger.Errorf("Token generation failed: %v", err)
		h.respondWithError(w, "Login failed", http.StatusInternalServerError)
		return
	}

	// Сохранение refresh-токена
	if err := h.repository.SaveRefreshToken(ctx, user.ID, tokens.RefreshToken); err != nil {
		h.logger.Errorf("Failed to save refresh token: %v", err)
		h.respondWithError(w, "Login failed", http.StatusInternalServerError)
		return
	}

	// Ответ
	response := LoginResponse{
		UserID: user.ID.String(),
		TokenPair: TokenPair{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}
	h.respondWithJSON(w, response, http.StatusOK)
}
