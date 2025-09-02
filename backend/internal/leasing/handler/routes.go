package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes возвращает новый роутер для обработчиков аренды.
func (h *LeasingHandler) Routes() http.Handler {
	r := chi.NewRouter()

	// POST / - создает новую аренду
	r.Post("/", h.CreateLease)
	// GET / - получает аренды текущего пользователя
	r.Get("/", h.GetMyLeases)

	return r
}
