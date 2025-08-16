package handler

import (
	"net/http"

	"github.com/rendley/vegshare/backend/pkg/api"
)

// --- Crop Handlers ---

func (h *FarmHandler) GetAllCrops(w http.ResponseWriter, r *http.Request) {
	crops, err := h.service.GetAllCrops(r.Context())
	if err != nil {
		h.logger.Errorf("ошибка при получении списка культур: %v", err)
		api.RespondWithError(w, "could not retrieve crops", http.StatusInternalServerError)
		return
	}
	api.RespondWithJSON(h.logger, w, crops, http.StatusOK)
}
