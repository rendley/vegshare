// Пакет handler содержит логику маршрутизации HTTP-запросов.
package handler

import (
	"database/sql"
	"github.com/rendley/backend/internal/auth/repository"
	"github.com/rendley/backend/pkg/security"
	"github.com/sirupsen/logrus"
)

// Handler — корневая структура для всех обработчиков.
// Позже добавим сюда зависимости (БД, логгер и т.д.)
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

//// SetupRoutes регистрирует все роуты.
//func (h *Handler) SetupRoutes(mux *http.ServeMux) {
//	mux.HandleFunc("GET /", h.homeHandler)
//	mux.HandleFunc("POST /register", h.registerHandler)
//	mux.HandleFunc("POST /login", h.loginHandler)
//}
