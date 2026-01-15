package middleware

import (
	"context"
	"net/http"
	"strings"

	"cinema-booking-system/internal/service"

	"go.uber.org/zap"
)

// contextKey is a custom type for context keys
type contextKey string

const (
	// UserContextKey is the key for storing user in context
	UserContextKey contextKey = "user"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	authService *service.AuthService
	logger      *zap.Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authService *service.AuthService, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
	}
}

// Authenticate verifies the JWT token and adds user to context
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.respondWithError(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		// Extract token (format: "Bearer <token>")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.respondWithError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		token := parts[1]

		// Validate token
		user, err := m.authService.ValidateToken(r.Context(), token)
		if err != nil {
			m.logger.Warn("Invalid token", zap.Error(err))
			m.respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), UserContextKey, user)

		// Call next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// respondWithError sends an error response
func (m *AuthMiddleware) respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(`{"success": false, "error": "` + message + `"}`))
}
