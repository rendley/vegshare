package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes returns a new router for the operations handler.
func (h *OperationsHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Route("/plots/{plotID}", func(r chi.Router) {
		r.Post("/plantings", h.PlantCrop)
		r.Get("/plantings", h.GetPlotCrops)
		r.Delete("/plantings/{plantingID}", h.RemoveCrop)
		r.Post("/actions", h.PerformAction)
	})

	return r
}