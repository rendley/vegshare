package handler

import "github.com/go-chi/chi/v5"

// RegisterRoutes registers all endpoints for the catalog handler.
func (h *CatalogHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/crops", func(r chi.Router) {
		r.Get("/", h.GetAllCrops)
		r.Post("/", h.CreateCrop)
	})
}
