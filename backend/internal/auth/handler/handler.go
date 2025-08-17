package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rendley/vegshare/backend/internal/auth/models"
	"github.com/rendley/vegshare/backend/internal/auth/service"
	"github.com/rendley/vegshare/backend/pkg/api"
	"github.com/sirupsen/logrus"
)

// --- DTOs ---

// RegisterRequest defines the request body for user registration.
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest defines the request body for user login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// --- Handler ---

// AuthHandler handles authentication requests.
type AuthHandler struct {
	service  service.AuthService
	logger   *logrus.Logger
	validate *validator.Validate
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(service service.AuthService, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		service:  service,
		logger:   logger,
		validate: validator.New(),
	}
}

// Register handles user registration.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, tokens, err := h.service.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		h.logger.Errorf("failed to register user: %v", err)
		if errors.Is(err, models.ErrUserExists) {
			api.RespondWithError(w, "user with this email already exists", http.StatusConflict)
			return
		}
		api.RespondWithError(w, "failed to register user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"user_id":       createdUser.ID,
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	}
	api.RespondWithJSON(h.logger, w, response, http.StatusCreated)
}

// Login handles user login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	loggedInUser, tokens, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Errorf("failed to login user: %v", err)
		if errors.Is(err, models.ErrInvalidCredentials) {
			api.RespondWithError(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		api.RespondWithError(w, "failed to login", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"user_id":       loggedInUser.ID,
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	}
	api.RespondWithJSON(h.logger, w, response, http.StatusOK)
}