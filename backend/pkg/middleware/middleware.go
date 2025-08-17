package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/pkg/config"
)

// contextKey is a custom type to avoid key collisions in context.
type contextKey string

const UserIDKey contextKey = "userID"

type Middleware struct {
	cfg *config.Config
}

func NewMiddleware(cfg *config.Config) *Middleware {
	return &Middleware{cfg: cfg}
}

// AuthMiddleware is a middleware to authenticate users using JWT.
func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString := headerParts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userIDStr, ok := claims["sub"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}