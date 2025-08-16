// Пакет handler
package handler

import "github.com/go-chi/chi/v5"

// RegisterRouter регистрирует все эндпоинты для обработчика фермы.
// Он принимает в качестве аргумента chi.Router, к которому привязывает пути.
func (h *FarmHandler) RegisterRouter(r chi.Router) {
	// Группируем маршруты, связанные с фермами
	r.Route("/api/v1/farms", func(r chi.Router) {
		r.Post("/", h.CreateFarm)   // POST /api/v1/farms
		r.Get("/", h.GetAllFarms) // GET /api/v1/farms
	})
}
