package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/user/service"
	"github.com/rendley/vegshare/backend/pkg/api"
	"github.com/rendley/vegshare/backend/pkg/middleware"
	"github.com/sirupsen/logrus"
)

// --- DTO (Data Transfer Objects) ---

type UpdateProfileRequest struct {
	Name      string `json:"name" validate:"omitempty,min=2,max=100"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin user"`
}

type ProfileResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// --- Handler ---

type UserHandler struct {
	service  service.UserService
	logger   *logrus.Logger
	validate *validator.Validate
}

func NewUserHandler(service service.UserService, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		service:  service,
		logger:   logger,
		validate: validator.New(),
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		api.RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userProfile, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("ошибка при получении пользователя: %v", err)
		api.RespondWithError(w, "User not found", http.StatusNotFound)
		return
	}

	api.RespondWithJSON(h.logger, w, userProfile, http.StatusOK)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		api.RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedUser, err := h.service.UpdateUser(r.Context(), userID, req.Name, req.AvatarURL)
	if err != nil {
		h.logger.Errorf("ошибка при обновлении пользователя: %v", err)
		api.RespondWithError(w, "Could not update user profile", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, updatedUser, http.StatusOK)
}

func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		api.RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.service.DeleteUser(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("ошибка при удалении пользователя: %v", err)
		api.RespondWithError(w, "Could not delete user account", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, map[string]string{"message": "User account deleted successfully"}, http.StatusOK)
}

// --- Admin Handlers ---

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers(r.Context())
	if err != nil {
		h.logger.Errorf("ошибка при получении списка пользователей: %v", err)
		api.RespondWithError(w, "Could not retrieve users", http.StatusInternalServerError)
		return
	}
	api.RespondWithJSON(h.logger, w, users, http.StatusOK)
}

func (h *UserHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid user ID in URL", http.StatusBadRequest)
		return
	}

	var req UpdateUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedUser, err := h.service.UpdateUserRole(r.Context(), userID, req.Role)
	if err != nil {
		h.logger.Errorf("ошибка при обновлении роли пользователя: %v", err)
		api.RespondWithError(w, "could not update user role", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, updatedUser, http.StatusOK)
}
