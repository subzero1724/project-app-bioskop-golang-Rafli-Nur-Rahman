package handler

import (
	"net/http"

	"cinema-booking-system/internal/service"
	"cinema-booking-system/internal/utils"

	"go.uber.org/zap"
)

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler struct {
	paymentService *service.PaymentService
	logger         *zap.Logger
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentService *service.PaymentService, logger *zap.Logger) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		logger:         logger,
	}
}

// GetAllPaymentMethods retrieves all available payment methods
// GET /api/payment-methods
func (h *PaymentHandler) GetAllPaymentMethods(w http.ResponseWriter, r *http.Request) {
	methods, err := h.paymentService.GetAllPaymentMethods(r.Context())
	if err != nil {
		h.logger.Error("Failed to get payment methods", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get payment methods")
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, methods, "")
}
