package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/pkg/api"
)

// --- Plot Handlers ---

func (h *FarmHandler) CreatePlotForGreenhouse(w http.ResponseWriter, r *http.Request) {
	greenhouseIDStr := chi.URLParam(r, "greenhouseID")
	greenhouseID, err := uuid.Parse(greenhouseIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid greenhouse ID in URL", http.StatusBadRequest)
		return
	}

	var req PlotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	plot, err := h.service.CreatePlot(r.Context(), req.Name, req.Size, greenhouseID)
	if err != nil {
		h.logger.Errorf("ошибка при создании грядки: %v", err)
		api.RespondWithError(w, "could not create plot", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, plot, http.StatusCreated)
}

func (h *FarmHandler) GetPlotsByGreenhouse(w http.ResponseWriter, r *http.Request) {
	greenhouseIDStr := chi.URLParam(r, "greenhouseID")
	greenhouseID, err := uuid.Parse(greenhouseIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid greenhouse ID in URL", http.StatusBadRequest)
		return
	}

	plots, err := h.service.GetPlotsByGreenhouse(r.Context(), greenhouseID)
	if err != nil {
		h.logger.Errorf("ошибка при получении списка грядок: %v", err)
		api.RespondWithError(w, "could not retrieve plots", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, plots, http.StatusOK)
}

func (h *FarmHandler) GetPlotByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "plotID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid plot ID in URL", http.StatusBadRequest)
		return
	}

	plot, err := h.service.GetPlotByID(r.Context(), id)
	if err != nil {
		h.logger.Errorf("ошибка при получении грядки: %v", err)
		api.RespondWithError(w, "plot not found", http.StatusNotFound)
		return
	}

	api.RespondWithJSON(h.logger, w, plot, http.StatusOK)
}

func (h *FarmHandler) UpdatePlot(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "plotID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid plot ID in URL", http.StatusBadRequest)
		return
	}

	var req UpdatePlotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	plot, err := h.service.UpdatePlot(r.Context(), id, req.Name, req.Size, req.Status)
	if err != nil {
		h.logger.Errorf("ошибка при обновлении грядки: %v", err)
		api.RespondWithError(w, "could not update plot", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, plot, http.StatusOK)
}

func (h *FarmHandler) DeletePlot(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "plotID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid plot ID in URL", http.StatusBadRequest)
		return
	}

	err = h.service.DeletePlot(r.Context(), id)
	if err != nil {
		h.logger.Errorf("ошибка при удалении грядки: %v", err)
		api.RespondWithError(w, "could not delete plot", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
