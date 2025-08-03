package handler

import (
	"encoding/json"
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
func (h *Handler) respondWithJSON(w http.ResponseWriter, statuCode int, payload interface{}) {
	// 1. Устанавливаем заголовки
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuCode)

	// 2. Кодируем и отправляем данные
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.logger.Errorf("Failed to encode JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
