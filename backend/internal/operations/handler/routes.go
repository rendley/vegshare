package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes returns a new router for the operations handler.
func (h *OperationsHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/actions", h.CreateAction)
	r.Get("/units/{unitID}/actions", h.GetActionsForUnit)
	r.Delete("/actions/{actionID}", h.CancelAction)

	return r
}
