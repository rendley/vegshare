package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/pkg/api"
)

func (h *OperationsHandler) GetPlotCrops(w http.ResponseWriter, r *http.Request) {
	plotIDStr := chi.URLParam(r, "plotID")
	plotID, err := uuid.Parse(plotIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid plot ID in URL", http.StatusBadRequest)
		return
	}

	plotCrops, err := h.service.GetPlotCrops(r.Context(), plotID)
	if err != nil {
		h.logger.Errorf("ошибка при получении посадок: %v", err)
		api.RespondWithError(w, "could not retrieve plot crops", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, plotCrops, http.StatusOK)
}
