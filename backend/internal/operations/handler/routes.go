package handler

import "github.com/go-chi/chi/v5"

// RegisterRouter регистрирует все эндпоинты для обработчика.
func (h *OperationsHandler) RegisterRouter(r chi.Router) {
	r.Post("/api/v1/plots/{plotID}/plantings", h.PlantCrop)
}