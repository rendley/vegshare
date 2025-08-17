package handler

import "github.com/go-chi/chi/v5"

// RegisterRouter регистрирует все эндпоинты для обработчика аренды.
func (h *LeasingHandler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/plots/{plotID}/lease", h.LeasePlot)
	r.Get("/api/v1/me/leases", h.GetMyLeases)
}