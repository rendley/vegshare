package handler

import "github.com/go-chi/chi/v5"

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/users/me", func(r chi.Router) {
		// Здесь в будущем можно будет добавить middleware для аутентификации
		r.Get("/", h.GetProfile)
		r.Patch("/", h.UpdateProfile)
		r.Delete("/", h.DeleteAccount)
	})
}
