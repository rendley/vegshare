// Пакет handler
package handler

import "github.com/go-chi/chi/v5"

// RegisterRouter регистрирует все эндпоинты для обработчика фермы.
func (h *FarmHandler) RegisterRouter(r chi.Router) {
	// Группируем маршруты, связанные с культурами
	r.Route("/api/v1/crops", func(r chi.Router) {
		r.Get("/", h.GetAllCrops) // GET /api/v1/crops
	})

	// Группируем маршруты, связанные с регионами
	r.Route("/api/v1/regions", func(r chi.Router) {
		r.Post("/", h.CreateRegion)       // POST /api/v1/regions
		r.Get("/", h.GetAllRegions)      // GET /api/v1/regions
		r.Get("/{id}", h.GetRegionByID)   // GET /api/v1/regions/{id}
		r.Put("/{id}", h.UpdateRegion)    // PUT /api/v1/regions/{id}
		r.Delete("/{id}", h.DeleteRegion) // DELETE /api/v1/regions/{id}
	})
}
