package service

import (
	"context"
	"fmt"

	"cinema-booking-system/internal/models"
	"cinema-booking-system/internal/repository"

	"go.uber.org/zap"
)

// PaymentService handles payment-related business logic
type PaymentService struct {
	paymentRepo *repository.PaymentRepository
	logger      *zap.Logger
}

// NewPaymentService creates a new payment service
func NewPaymentService(paymentRepo *repository.PaymentRepository, logger *zap.Logger) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		logger:      logger,
	}
}

// GetAllPaymentMethods retrieves all available payment methods
func (s *PaymentService) GetAllPaymentMethods(ctx context.Context) ([]*models.PaymentMethod, error) {
	methods, err := s.paymentRepo.GetAllPaymentMethods(ctx)
	if err != nil {
		s.logger.Error("Failed to get payment methods", zap.Error(err))
		return nil, fmt.Errorf("failed to get payment methods")
	}

	return methods, nil
}
