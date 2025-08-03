package handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

// respondWithError - универсальный метод для отправки ошибок
func (h *Handler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	// 1. Логируем ошибку
	h.logger.Warnf("Error %d: %s", statusCode, message)

	// 2. Устанавливаем заголовки
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// 3. Формируем и отправляем JSON
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// respondWithJSON - универсальный метод для успешных ответов
func (h *Handler) respondWithJSON(w http.ResponseWriter, payload interface{}, statusCode int) {
	// 1. Устанавливаем заголовки
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// 2. Кодируем и отправляем данные
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.logger.Errorf("Failed to encode JSON response: %v", err)
	}
}

func (h *Handler) generateTokens(userID uuid.UUID) (*TokenPair, error) {
	// Генерируем access token
	accessToken, err := h.jwtGenerator.GenerateAccessToken(userID)
	if err != nil {
		h.logger.Errorf("Failed to generate access token: %v", err)
		return nil, fmt.Errorf("token generation failed")
	}

	// Генерируем refresh token
	refreshToken, err := h.jwtGenerator.GenerateRefreshToken()
	if err != nil {
		h.logger.Errorf("Failed to generate refresh token: %v", err)
		return nil, fmt.Errorf("token generation failed")
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
