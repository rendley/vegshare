package handler

import "github.com/go-chi/chi/v5"

// RegisterRouter регистрирует все эндпоинты для обработчика аренды.
func (h *LeasingHandler) RegisterRouter(r chi.Router) {
	// r.Post("/api/v1/plots/{plotID}/lease", h.LeasePlot) // <-- добавим в следующей задаче
}
