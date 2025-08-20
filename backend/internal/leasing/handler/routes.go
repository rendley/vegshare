package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes returns a new router for the leasing handler.
func (h *LeasingHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/plots/{plotID}/lease", h.LeasePlot)
	r.Get("/me/leases", h.GetMyLeases)

	return r
}