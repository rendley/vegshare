package handler

import "github.com/go-chi/chi/v5"

// RegisterRoutes регистрирует все эндпоинты для обработчика фермы.
func (h *FarmHandler) RegisterRoutes(r chi.Router) {
	

	// --- Маршруты для Регионов ---
	r.Route("/api/v1/regions", func(r chi.Router) {
		r.Post("/", h.CreateRegion)
		r.Get("/", h.GetAllRegions)

		// Маршруты для конкретного региона
		r.Route("/{regionID}", func(r chi.Router) {
			r.Get("/", h.GetRegionByID)
			r.Put("/", h.UpdateRegion)
			r.Delete("/", h.DeleteRegion)

			// Вложенные маршруты для земельных участков
			r.Get("/land-parcels", h.GetLandParcelsByRegion)
			r.Post("/land-parcels", h.CreateLandParcelForRegion)
		})
	})

	// --- Маршруты для Земельных участков ---
	r.Route("/api/v1/land-parcels", func(r chi.Router) {
		// Маршруты для конкретного земельного участка
		r.Route("/{parcelID}", func(r chi.Router) {
			r.Get("/", h.GetLandParcelByID)
			r.Put("/", h.UpdateLandParcel)
			r.Delete("/", h.DeleteLandParcel)

			// Вложенные маршруты для теплиц
			r.Get("/greenhouses", h.GetGreenhousesByLandParcel)
			r.Post("/greenhouses", h.CreateGreenhouseForLandParcel)
		})
	})

	// --- Маршруты для Теплиц ---
	r.Route("/api/v1/greenhouses", func(r chi.Router) {
		// Маршруты для конкретной теплицы
		r.Route("/{greenhouseID}", func(r chi.Router) {
			r.Get("/", h.GetGreenhouseByID)
			r.Put("/", h.UpdateGreenhouse)
			r.Delete("/", h.DeleteGreenhouse)

			// Вложенные маршруты для грядок
			r.Get("/plots", h.GetPlotsByGreenhouse)
			r.Post("/plots", h.CreatePlotForGreenhouse)
		})
	})

	// --- Маршруты для Грядок ---
	r.Route("/api/v1/plots", func(r chi.Router) {
		r.Route("/{plotID}", func(r chi.Router) {
			r.Get("/", h.GetPlotByID)
			r.Put("/", h.UpdatePlot)
			r.Delete("/", h.DeletePlot)
		})
	})
}