// Пакет handler содержит логику маршрутизации HTTP-запросов.
package handler

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Handler — корневая структура для всех обработчиков.
// Позже добавим сюда зависимости (БД, логгер и т.д.)
type Handler struct {
	db     *sql.DB        // Подключение к PostgreSQL
	logger *logrus.Logger // Логгер
}

// New создаёт экземпляр Handler c зависимостями.
func New(db *sql.DB, logger *logrus.Logger) *Handler {
	return &Handler{
		db:     db,
		logger: logger,
	}
}

// SetupRoutes регистрирует все роуты.
func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.homeHandler)
	mux.HandleFunc("POST /register", h.registerHandler)
	mux.HandleFunc("POST /login", h.loginHandler)
}
