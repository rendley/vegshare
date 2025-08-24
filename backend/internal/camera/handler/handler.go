package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/camera/service"
	"github.com/rendley/vegshare/backend/pkg/api"
	"github.com/sirupsen/logrus"
)

// --- DTOs ---
type CreateCameraRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=100"`
	RTSPURL string `json:"rtsp_url" validate:"required,url"`
}

// --- Handler ---

type CameraHandler struct {
	service  service.Service
	logger   *logrus.Logger
	validate *validator.Validate
}

func NewCameraHandler(s service.Service, l *logrus.Logger) *CameraHandler {
	return &CameraHandler{
		service:  s,
		logger:   l,
		validate: validator.New(),
	}
}

func (h *CameraHandler) CreateCamera(w http.ResponseWriter, r *http.Request) {
	plotIDStr := chi.URLParam(r, "plotID")
	plotID, err := uuid.Parse(plotIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid plot ID in URL", http.StatusBadRequest)
		return
	}

	var req CreateCameraRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	camera, err := h.service.CreateCamera(r.Context(), req.Name, req.RTSPURL, plotID)
	if err != nil {
		h.logger.Errorf("ошибка при создании камеры: %v", err)
		api.RespondWithError(w, "could not create camera", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, camera, http.StatusCreated)
}

func (h *CameraHandler) GetCamerasByPlotID(w http.ResponseWriter, r *http.Request) {
	plotIDStr := chi.URLParam(r, "plotID")
	plotID, err := uuid.Parse(plotIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid plot ID in URL", http.StatusBadRequest)
		return
	}

	cameras, err := h.service.GetCamerasByPlotID(r.Context(), plotID)
	if err != nil {
		h.logger.Errorf("ошибка при получении списка камер: %v", err)
		api.RespondWithError(w, "could not retrieve cameras", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, cameras, http.StatusOK)
}

func (h *CameraHandler) DeleteCamera(w http.ResponseWriter, r *http.Request) {
	cameraIDStr := chi.URLParam(r, "cameraID")
	cameraID, err := uuid.Parse(cameraIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid camera ID in URL", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteCamera(r.Context(), cameraID)
	if err != nil {
		h.logger.Errorf("ошибка при удалении камеры: %v", err)
		api.RespondWithError(w, "could not delete camera", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
