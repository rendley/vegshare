// Пакет handler содержит логику маршрутизации HTTP-запросов.
package handler

import (
	"net/http"
)

// Handler — корневая структура для всех обработчиков.
// Позже добавим сюда зависимости (БД, логгер и т.д.)
type Handler struct{}

// New создаёт экземпляр Handler.
func New() *Handler {
	return &Handler{}
}

// SetupRoutes регистрирует все роуты.
func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.homeHandler)
	mux.HandleFunc("POST /register", h.registerHandler)
	mux.HandleFunc("POST /login", h.loginHandler)
}
