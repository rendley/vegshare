package jwt

import "time"

// Generator отвечает за создание и проверку токенов
type Generator struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTl time.Duration
}

// NewGenerator создает новый экземпляр генератора токенов
func NewGenerator(secret string, accessTTL, refreshTTL time.Duration) *Generator {
	return &Generator{
		secret:          []byte(secret),
		accessTokenTTL:  accessTTL,
		refreshTokenTTl: refreshTTL,
	}
}

// TokenPair содержит оба типа токенов
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
