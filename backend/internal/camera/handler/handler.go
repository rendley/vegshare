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
// CreateCameraRequest определяет тело запроса для создания новой камеры.
type CreateCameraRequest struct {
	Name         string    `json:"name" validate:"required,min=2,max=100"`
	RTSPPathName string    `json:"rtsp_path_name" validate:"required,min=1,max=100"`
	UnitID       uuid.UUID `json:"unit_id" validate:"required"`
	UnitType     string    `json:"unit_type" validate:"required"`
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

// CreateCamera обрабатывает POST-запросы для создания камеры для юнита.
func (h *CameraHandler) CreateCamera(w http.ResponseWriter, r *http.Request) {
	var req CreateCameraRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	camera, err := h.service.CreateCamera(r.Context(), req.Name, req.RTSPPathName, req.UnitID, req.UnitType)
	if err != nil {
		h.logger.Errorf("ошибка при создании камеры: %v", err)
		api.RespondWithError(w, "could not create camera", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, camera, http.StatusCreated)
}

// GetCameras обрабатывает GET-запросы для получения списка камер по ID и типу юнита.
func (h *CameraHandler) GetCameras(w http.ResponseWriter, r *http.Request) {
	unitIDStr := r.URL.Query().Get("unit_id")
	unitType := r.URL.Query().Get("unit_type")

	if unitIDStr == "" || unitType == "" {
		api.RespondWithError(w, "query parameters 'unit_id' and 'unit_type' are required", http.StatusBadRequest)
		return
	}

	unitID, err := uuid.Parse(unitIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid unit_id format", http.StatusBadRequest)
		return
	}

	cameras, err := h.service.GetCamerasByUnitID(r.Context(), unitID, unitType)
	if err != nil {
		h.logger.Errorf("ошибка при получении списка камер: %v", err)
		api.RespondWithError(w, "could not retrieve cameras", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, cameras, http.StatusOK)
}

// DeleteCamera обрабатывает DELETE-запросы для удаления камеры по ее ID.
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