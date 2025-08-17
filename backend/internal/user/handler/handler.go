package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/user/service"
	"github.com/rendley/vegshare/backend/pkg/api"
	"github.com/sirupsen/logrus"
)

type contextKey string

const userIDKey contextKey = "userID"

// --- DTO (Data Transfer Objects) ---

type UpdateProfileRequest struct {
	FullName  string `json:"full_name" validate:"omitempty,min=2,max=100"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
}

type ProfileResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
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
	userID, _ := uuid.Parse("aec44274-55b0-497a-a18b-7f8efe7d8a9e")

	userProfile, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("ошибка при получении пользователя: %v", err)
		api.RespondWithError(w, "User not found", http.StatusNotFound)
		return
	}

	api.RespondWithJSON(h.logger, w, userProfile, http.StatusOK)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
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

	updatedUser, err := h.service.UpdateUser(r.Context(), userID, req.FullName, req.AvatarURL)
	if err != nil {
		h.logger.Errorf("ошибка при обновлении пользователя: %v", err)
		api.RespondWithError(w, "Could not update user profile", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, updatedUser, http.StatusOK)
}

func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
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
