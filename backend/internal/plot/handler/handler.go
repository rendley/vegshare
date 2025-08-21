package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/plot/service"
	"github.com/rendley/vegshare/backend/pkg/api"
	"github.com/sirupsen/logrus"
)

// --- DTOs ---
type PlotRequest struct {
	Name         string    `json:"name" validate:"required,min=2,max=100"`
	Size         string    `json:"size" validate:"omitempty,min=1,max=50"`
	GreenhouseID uuid.UUID `json:"greenhouse_id" validate:"required"`
}

type UpdatePlotRequest struct {
	Name   string `json:"name" validate:"required,min=2,max=100"`
	Size   string `json:"size" validate:"omitempty,min=1,max=50"`
	Status string `json:"status" validate:"required,oneof=available rented maintenance"`
}

// --- Handler ---

type PlotHandler struct {
	service  service.Service
	logger   *logrus.Logger
	validate *validator.Validate
}

func NewPlotHandler(s service.Service, l *logrus.Logger) *PlotHandler {
	return &PlotHandler{
		service:  s,
		logger:   l,
		validate: validator.New(),
	}
}

func (h *PlotHandler) CreatePlot(w http.ResponseWriter, r *http.Request) {
	var req PlotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	plot, err := h.service.CreatePlot(r.Context(), req.Name, req.Size, req.GreenhouseID)
	if err != nil {
		h.logger.Errorf("ошибка при создании грядки: %v", err)
		api.RespondWithError(w, "could not create plot", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, plot, http.StatusCreated)
}

func (h *PlotHandler) GetPlots(w http.ResponseWriter, r *http.Request) {
	greenhouseIDStr := r.URL.Query().Get("greenhouse_id")
	if greenhouseIDStr == "" {
		// Handle case where we might want to list all plots in the future
		api.RespondWithError(w, "greenhouse_id query parameter is required", http.StatusBadRequest)
		return
	}

	greenhouseID, err := uuid.Parse(greenhouseIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid greenhouse_id query parameter", http.StatusBadRequest)
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

func (h *PlotHandler) GetPlotByID(w http.ResponseWriter, r *http.Request) {
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

func (h *PlotHandler) UpdatePlot(w http.ResponseWriter, r *http.Request) {
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

func (h *PlotHandler) DeletePlot(w http.ResponseWriter, r *http.Request) {
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
