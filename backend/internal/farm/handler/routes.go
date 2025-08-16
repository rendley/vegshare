// Пакет handler
package handler

import "github.com/go-chi/chi/v5"

// RegisterRouter регистрирует все эндпоинты для обработчика фермы.
func (h *FarmHandler) RegisterRouter(r chi.Router) {
	// --- Маршруты для Культур ---
	r.Route("/api/v1/crops", func(r chi.Router) {
		r.Get("/", h.GetAllCrops)
	})

	// --- Маршруты для Регионов и вложенных Земельных участков ---
	r.Route("/api/v1/regions", func(r chi.Router) {
		r.Post("/", h.CreateRegion)
		r.Get("/", h.GetAllRegions)

		// Маршруты, специфичные для одного региона, включая вложенные
		r.Route("/{regionID}", func(r chi.Router) {
			r.Get("/", h.GetRegionByID)
			r.Put("/", h.UpdateRegion)
			r.Delete("/", h.DeleteRegion)

			// Вложенные маршруты для земельных участков
			// POST /api/v1/regions/{regionID}/land-parcels
			r.Post("/land-parcels", h.CreateLandParcel)
			// GET /api/v1/regions/{regionID}/land-parcels
			r.Get("/land-parcels", h.GetLandParcelsByRegion)
		})
	})

	// --- Маршруты для прямого управления Земельными участками по их ID ---
	r.Route("/api/v1/land-parcels", func(r chi.Router) {
		r.Get("/{id}", h.GetLandParcelByID)
		r.Put("/{id}", h.UpdateLandParcel)
		r.Delete("/{id}", h.DeleteLandParcel)
	})
}
