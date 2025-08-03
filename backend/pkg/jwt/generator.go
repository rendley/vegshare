package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type Generator interface {
	GenerateAccessToken(userID uuid.UUID) (string, error)
	GenerateRefreshToken() (string, error)
}

type JWTGenerator struct {
	secretKey     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewGenerator(secretKey string, accessExpiry, refreshExpiry time.Duration) *JWTGenerator {
	return &JWTGenerator{
		secretKey:     secretKey,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

func (g *JWTGenerator) GenerateAccessToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(g.accessExpiry).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(g.secretKey))
}

func (g *JWTGenerator) GenerateRefreshToken() (string, error) {
	token := uuid.New().String() + uuid.New().String() // Пример генерации
	return token, nil
}
