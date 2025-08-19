
package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes returns a new router for the catalog handler.
func (h *CatalogHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/crops", h.GetAllCrops)
	r.Post("/crops", h.CreateCrop)

	return r
}

