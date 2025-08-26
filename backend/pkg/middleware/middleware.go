package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/pkg/config"
	"github.com/sirupsen/logrus"
)

// ContextKey is a custom type to avoid key collisions in context.
type ContextKey string

const (
	// UserIDKey is the key for the user ID in the context.
	UserIDKey     ContextKey = "userID"
	// UserClaimsKey is the key for the user claims in the context.
	UserClaimsKey ContextKey = "userClaims"
)

type Middleware struct {
	cfg    *config.Config
	logger *logrus.Logger
}

func NewMiddleware(cfg *config.Config, logger *logrus.Logger) *Middleware {
	return &Middleware{cfg: cfg, logger: logger}
}

// AuthMiddleware is a middleware to authenticate users using JWT.
func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.logger.Info("Auth header not found")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			m.logger.Info("Invalid auth header format")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString := headerParts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			m.logger.Info("Invalid token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			m.logger.Info("Invalid token claims")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userIDStr, ok := claims["sub"].(string)
		if !ok {
			m.logger.Info("Invalid userID in token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			m.logger.Info("Failed to parse userID")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		m.logger.Infof("userID from token: %s", userID)

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminMiddleware checks if the user has admin role.
// This middleware must be used AFTER AuthMiddleware.
func (m *Middleware) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(UserClaimsKey).(jwt.MapClaims)
		if !ok {
			m.logger.Error("Admin middleware error: claims not found in context. Is AuthMiddleware running before?")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			m.logger.Warnf("Forbidden: user %s with role '%s' tried to access admin route", claims["sub"], role)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CorsMiddleware добавляет заголовки CORS к ответам.
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Если это pre-flight запрос (OPTIONS), просто возвращаем OK.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// QueryParamAuthMiddleware is a middleware to authenticate users by checking the "token" query parameter.
func (m *Middleware) QueryParamAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.URL.Query().Get("token")

		if tokenString == "" {
			m.logger.Info("Auth token not found in query param")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			m.logger.Info("Invalid token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			m.logger.Info("Invalid token claims")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userIDStr, ok := claims["sub"].(string)
		if !ok {
			m.logger.Info("Invalid userID in token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			m.logger.Info("Failed to parse userID")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		m.logger.Infof("userID from token: %s", userID)

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext извлекает ID пользователя из контекста запроса.
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}
