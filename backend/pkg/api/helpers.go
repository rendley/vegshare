package api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

// RespondWithError - универсальный метод для отправки ошибок
func RespondWithError(w http.ResponseWriter, message string, statusCode int) {
	//1. Устанавливаем заголовки
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	//2. Формируем и отправляем JSON
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// RespondWithJSON - универсальный метод для успешных ответов
func RespondWithJSON(logger *logrus.Logger, w http.ResponseWriter, payload interface{}, statusCode int) {
	// 1. Устанавливаем заголовки
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// 2. Kодируем и отправляем данные
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logger.Errorf("Failed to encode JSON reponse: %v", err)
	}

}

//// TokenGenerator интерфейс для генерации токенов
//type TokenGenerator interface {
//	GenerateAccessToken(userID uuid.UUID) (string, error)
//	GenerateRefreshToken() (string, error)
//}
//
//// GenerateTokens - универсальный метод генерации токенов
//func GenerateTokens(logger *logrus.Logger, generator TokenGenerator, userID uuid.UUID) (*TokenPair, error) {
//	// Генерируем access token
//	accessToken, err := generator.GenerateAccessToken(userID)
//	if err != nil {
//		logger.Errorf("Failed to generate access token: %v", err)
//		return nil, fmt.Errorf("token generation failed")
//	}
//
//	// Генерируем refresh token
//	refreshToken, err := generator.GenerateRefreshToken()
//	if err != nil {
//		logger.Errorf("Failed to generate refresh token: %v", err)
//		return nil, fmt.Errorf("token generation failed")
//	}
//
//	return &TokenPair{
//		AccessToken:  accessToken,
//		RefreshToken: refreshToken,
//	}, nil
//}
//
//type TokenPair struct {
//	AccessToken  string `json:"access_token"`
//	RefreshToken string `json:"refresh_token"`
//}
