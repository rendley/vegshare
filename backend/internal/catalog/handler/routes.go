package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes returns a new router for the catalog handler.
func (h *CatalogHandler) Routes() http.Handler {
	r := chi.NewRouter()

	// GET /items?type=crop - получить все элементы по типу
	r.Get("/items", h.GetItems)

	// POST /items - создать новый элемент
	r.Post("/items", h.CreateItem) // В будущем защитить middleware для админа

	return r
}