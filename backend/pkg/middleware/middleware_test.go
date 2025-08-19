package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/pkg/config"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a mock handler that checks the context for a user ID.
func newMockHandlerWithContextCheck(t *testing.T, expectedUserID uuid.UUID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
		assert.True(t, ok, "User ID should be in context")
		assert.Equal(t, expectedUserID, userID, "User ID in context should match expected ID")
		w.WriteHeader(http.StatusOK)
	}
}

func TestAuthMiddleware(t *testing.T) {
	// Setup
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
		},
	}
	mw := NewMiddleware(cfg)

	t.Run("should pass with valid token", func(t *testing.T) {
		// Arrange
		expectedUserID := uuid.New()
		claims := jwt.MapClaims{
			"sub": expectedUserID.String(),
			"exp": time.Now().Add(time.Hour * 1).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
		assert.NoError(t, err)

		mockHandler := newMockHandlerWithContextCheck(t, expectedUserID)
		authHandler := mw.AuthMiddleware(mockHandler)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rr := httptest.NewRecorder()

		// Act
		authHandler.ServeHTTP(rr, req)

		// Assert
		assert.Equal(t, http.StatusOK, rr.Code, "Status code should be OK")
	})

	t.Run("should fail with no token", func(t *testing.T) {
		// Arrange
		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("next handler should not be called")
		})
		authHandler := mw.AuthMiddleware(mockHandler)

		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		// Act
		authHandler.ServeHTTP(rr, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Status code should be Unauthorized")
	})

	t.Run("should fail with invalid token", func(t *testing.T) {
		// Arrange
		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("next handler should not be called")
		})
		authHandler := mw.AuthMiddleware(mockHandler)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rr := httptest.NewRecorder()

		// Act
		authHandler.ServeHTTP(rr, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Status code should be Unauthorized")
	})
}

func TestCorsMiddleware(t *testing.T) {
	t.Run("should add CORS headers for non-OPTIONS request", func(t *testing.T) {
		// Arrange
		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		corsHandler := CorsMiddleware(mockHandler)

		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		// Act
		corsHandler.ServeHTTP(rr, req)

		// Assert
		assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
	})

	t.Run("should return 200 OK for OPTIONS request", func(t *testing.T) {
		// Arrange
		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("next handler should not be called for OPTIONS request")
		})
		corsHandler := CorsMiddleware(mockHandler)

		req := httptest.NewRequest("OPTIONS", "/", nil)
		rr := httptest.NewRecorder()

		// Act
		corsHandler.ServeHTTP(rr, req)

		// Assert
		assert.Equal(t, http.StatusOK, rr.Code, "Status code should be OK for OPTIONS request")
	})
}
