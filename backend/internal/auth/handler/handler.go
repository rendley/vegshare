// Пакет handler содержит базовую структуру и конструктор.
// В данном пакете также собраны методы обработчиков связанные с auth
package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/rendley/backend/internal/auth/repository"
	"github.com/rendley/backend/pkg/security"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Handler — корневая структура для всех обработчиков.
type Handler struct {
	db             *sql.DB                    // Подключение к PostgreSQL
	logger         *logrus.Logger             // Логгер
	passwordHasher security.PasswordHasher    // Хеширование пассводра
	repository     *repository.AuthRepository // репозиторий для работы с базой
}

// New создаёт экземпляр Handler c зависимостями.
func New(db *sql.DB, hasher security.PasswordHasher, logger *logrus.Logger) *Handler {
	return &Handler{
		db:             db,
		logger:         logger,
		passwordHasher: hasher,
		repository:     repository.NewAuthRepository(db), // Инициализируем репозиторий
	}
}

//#################################################################//

// homeHandler обрабатывает запросы к корневому пути ("/").
func (h *Handler) homeHandler(w http.ResponseWriter, r *http.Request) {
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
func (h *Handler) registerHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Логируем начало обработки запроса
	h.logger.Info("Register request received")

	// 2. Проверяем метод запроса (должен быть POST)
	if r.Method != http.MethodPost {
		h.logger.Warn("Wrong method used for registration")
		h.respondWithError(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	// 3. Читаем тело запроса делаем это из локальной структуры
	// стоит начинать с них а если будут повторения всегда можно вынести в types.go
	// как это сделано в loginHandler
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 4. Парсим JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Warn("Failed to parse JSON: %v", err)
		h.respondWithError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// 5. Проверяем обязательные поля
	if request.Email == "" || request.Password == "" {
		h.respondWithError(w, "Email and passwoed are required", http.StatusBadRequest)
		return
	}

	// 6. Хешируем пароль перед сохранением в БД
	hashedPassword, err := h.passwordHasher.Hash(request.Password)
	if err != nil {
		h.logger.Errorf("Password hashing failed: %v", err)
		h.respondWithError(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	// 7. Сохраняем пользователя в БД
	userID, err := h.repository.CreateUser(request.Email, hashedPassword)
	if err != nil {
		h.logger.Errorf("Failed to create user: %v", err)
		h.respondWithError(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	// 8. Формируем ответ (с заглушкой для токена)
	response := struct {
		Token  string `json:"token"`
		UserID string `json:"user_id"`
	}{
		Token:  "fake-jwt-token-for-registration",
		UserID: userID,
	}

	// 9. Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created — стандартный статус для успешной регистрации
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// loginHandler обрабатывает POST /login.
func (h *Handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// fmt.Fprint(w, "Login endpoint (will be implemented later)")

	// 2. Читаем и парсим JSON-тело.
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}
	// 3. Валидация (минимальная).
	if req.Email == "" || len(req.Password) < 6 {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}
	// 4. Заглушка: проверка логина/пароля (позже заменим на БД).
	if req.Email != "test@example.com" || req.Password != "password" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 5. Формируем успешный ответ.
	response := LoginResponse{
		Token:  "fake-jwt-token",
		UserID: "123",
	}
	//6. Сериализуем ответ в JSON и отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
