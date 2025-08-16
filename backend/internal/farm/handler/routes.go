// Пакет handler
package handler

import "github.com/go-chi/chi/v5"

// RegisterRouter регистрирует все эндпоинты для обработчика фермы.
func (h *FarmHandler) RegisterRouter(r chi.Router) {
	// Группируем маршруты, связанные с культурами
	r.Route("/api/v1/crops", func(r chi.Router) {
		r.Get("/", h.GetAllCrops) // GET /api/v1/crops
	})
}