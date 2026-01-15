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

// BookingHandler handles booking-related HTTP requests
type BookingHandler struct {
	bookingService *service.BookingService
	validator      *utils.Validator
	logger         *zap.Logger
}

// NewBookingHandler creates a new booking handler
func NewBookingHandler(bookingService *service.BookingService, validator *utils.Validator, logger *zap.Logger) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
		validator:      validator,
		logger:         logger,
	}
}

// CreateBooking creates a new seat reservation
// POST /api/booking
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.BookingRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode booking request", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := h.validator.Validate(req); err != nil {
		utils.RespondWithValidationError(w, err)
		return
	}

	// Create booking
	booking, err := h.bookingService.CreateBooking(r.Context(), user.ID, &req)
	if err != nil {
		h.logger.Error("Failed to create booking", zap.Int("user_id", user.ID), zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithSuccess(w, http.StatusCreated, booking, "Booking created successfully")
}

// GetUserBookings retrieves booking history for the logged-in user
// GET /api/user/bookings
func (h *BookingHandler) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get user bookings
	bookings, err := h.bookingService.GetUserBookings(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("Failed to get user bookings", zap.Int("user_id", user.ID), zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get bookings")
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, bookings, "")
}

// ProcessPayment processes payment for a booking
// POST /api/pay
func (h *BookingHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.PaymentRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode payment request", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := h.validator.Validate(req); err != nil {
		utils.RespondWithValidationError(w, err)
		return
	}

	// Process payment
	booking, err := h.bookingService.ProcessPayment(r.Context(), user.ID, &req)
	if err != nil {
		h.logger.Error("Failed to process payment", zap.Int("user_id", user.ID), zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, booking, "Payment processed successfully")
}
