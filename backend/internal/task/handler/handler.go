package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/task/service"
	"github.com/rendley/vegshare/backend/pkg/api"
	"github.com/rendley/vegshare/backend/pkg/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
)

type TaskHandler struct {
	service  service.Service
	logger   *logrus.Logger
	validate *validator.Validate
}

func NewTaskHandler(s service.Service, l *logrus.Logger) *TaskHandler {
	return &TaskHandler{
		service:  s,
		logger:   l,
		validate: validator.New(),
	}
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetAllTasks(r.Context())
	if err != nil {
		h.logger.Errorf("ошибка при получении всех задач: %v", err)
		api.RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, tasks, http.StatusOK)
}

func (h *TaskHandler) AcceptTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		api.RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	taskIDStr := chi.URLParam(r, "taskID")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		api.RespondWithError(w, "некорректный ID задачи", http.StatusBadRequest)
		return
	}

	task, err := h.service.AcceptTask(r.Context(), taskID, userID)
	if err != nil {
		h.logger.Errorf("ошибка при принятии задачи в работу: %v", err)
		api.RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, task, http.StatusOK)
}

func (h *TaskHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		api.RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	taskIDStr := chi.URLParam(r, "taskID")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		api.RespondWithError(w, "некорректный ID задачи", http.StatusBadRequest)
		return
	}

	task, err := h.service.CompleteTask(r.Context(), taskID, userID)
	if err != nil {
		h.logger.Errorf("ошибка при завершении задачи: %v", err)
		api.RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, task, http.StatusOK)
}

func (h *TaskHandler) FailTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		api.RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	taskIDStr := chi.URLParam(r, "taskID")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		api.RespondWithError(w, "некорректный ID задачи", http.StatusBadRequest)
		return
	}

	task, err := h.service.FailTask(r.Context(), taskID, userID)
	if err != nil {
		h.logger.Errorf("ошибка при провале задачи: %v", err)
		api.RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, task, http.StatusOK)
}
