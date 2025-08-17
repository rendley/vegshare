package handler

import "github.com/go-chi/chi/v5"

// RegisterRouter регистрирует все эндпоинты для обработчика.
func (h *OperationsHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/plots/{plotID}", func(r chi.Router) {
		r.Post("/plantings", h.PlantCrop)
		r.Get("/plantings", h.GetPlotCrops)
		r.Delete("/plantings/{plantingID}", h.RemoveCrop)
		r.Post("/actions", h.PerformAction)
	})
}