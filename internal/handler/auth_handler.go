package handler

import (
	"encoding/json"
	"net/http"

	"cinema-booking-system/internal/dto"
	"cinema-booking-system/internal/middleware"
	"cinema-booking-system/internal/models"
	"cinema-booking-system/internal/service"
	"cinema-booking-system/internal/utils"

	"go.uber.org/zap"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService *service.AuthService
	validator   *utils.Validator
	logger      *zap.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService, validator *utils.Validator, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
		logger:      logger,
	}
}

// Register handles user registration
// POST /api/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode register request", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := h.validator.Validate(req); err != nil {
		utils.RespondWithValidationError(w, err)
		return
	}

	// Register user
	user, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to register user", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Prepare response
	userResponse := dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
	}

	utils.RespondWithSuccess(w, http.StatusCreated, userResponse, "User registered successfully")
}

// Login handles user login
// POST /api/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode login request", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := h.validator.Validate(req); err != nil {
		utils.RespondWithValidationError(w, err)
		return
	}

	// Login user
	loginResponse, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to login user", zap.Error(err))
		utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, loginResponse, "Login successful")
}

// Logout handles user logout
// POST /api/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	token := authHeader[7:] // Remove "Bearer " prefix

	// Logout user
	if err := h.authService.Logout(r.Context(), token); err != nil {
		h.logger.Error("Failed to logout user", zap.Int("user_id", user.ID), zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, nil, "Logout successful")
}
