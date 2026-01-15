package service

import (
	"context"
	"fmt"

	"cinema-booking-system/internal/dto"
	"cinema-booking-system/internal/models"
	"cinema-booking-system/internal/repository"

	"go.uber.org/zap"
)

// BookingService handles booking-related business logic
type BookingService struct {
	bookingRepo *repository.BookingRepository
	cinemaRepo  *repository.CinemaRepository
	paymentRepo *repository.PaymentRepository
	logger      *zap.Logger
}

// NewBookingService creates a new booking service
func NewBookingService(
	bookingRepo *repository.BookingRepository,
	cinemaRepo *repository.CinemaRepository,
	paymentRepo *repository.PaymentRepository,
	logger *zap.Logger,
) *BookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		cinemaRepo:  cinemaRepo,
		paymentRepo: paymentRepo,
		logger:      logger,
	}
}

// CreateBooking creates a new seat reservation
func (s *BookingService) CreateBooking(ctx context.Context, userID int, req *dto.BookingRequest) (*models.Booking, error) {
	// Validate cinema exists
	cinema, err := s.cinemaRepo.GetByID(ctx, req.CinemaID)
	if err != nil {
		s.logger.Error("Cinema not found", zap.Int("cinema_id", req.CinemaID))
		return nil, fmt.Errorf("cinema not found")
	}

	// Validate seat exists and belongs to cinema
	seat, err := s.cinemaRepo.GetSeatByID(ctx, req.SeatID)
	if err != nil {
		s.logger.Error("Seat not found", zap.Int("seat_id", req.SeatID))
		return nil, fmt.Errorf("seat not found")
	}

	if seat.CinemaID != cinema.ID {
		s.logger.Error("Seat does not belong to cinema",
			zap.Int("seat_id", req.SeatID),
			zap.Int("cinema_id", req.CinemaID))
		return nil, fmt.Errorf("seat does not belong to the specified cinema")
	}

	// Validate payment method
	isValid, err := s.paymentRepo.ValidatePaymentMethod(ctx, req.PaymentMethod)
	if err != nil || !isValid {
		s.logger.Error("Invalid payment method", zap.Int("payment_method_id", req.PaymentMethod))
		return nil, fmt.Errorf("invalid payment method")
	}

	// Check seat availability
	isAvailable, err := s.bookingRepo.CheckSeatAvailability(ctx, req.CinemaID, req.SeatID, req.Date, req.Time+":00")
	if err != nil {
		s.logger.Error("Failed to check seat availability", zap.Error(err))
		return nil, fmt.Errorf("failed to check seat availability")
	}

	if !isAvailable {
		s.logger.Warn("Seat already booked",
			zap.Int("cinema_id", req.CinemaID),
			zap.Int("seat_id", req.SeatID),
			zap.String("date", req.Date),
			zap.String("time", req.Time))
		return nil, fmt.Errorf("seat is already booked for the specified time")
	}

	// Create booking
	booking := &models.Booking{
		UserID:          userID,
		CinemaID:        req.CinemaID,
		SeatID:          req.SeatID,
		BookingDate:     req.Date,
		BookingTime:     req.Time + ":00",
		PaymentMethodID: &req.PaymentMethod,
		PaymentStatus:   "pending",
		TotalAmount:     seat.Price,
		BookingStatus:   "reserved",
	}

	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		s.logger.Error("Failed to create booking", zap.Error(err))
		return nil, fmt.Errorf("failed to create booking")
	}

	s.logger.Info("Booking created successfully",
		zap.Int("booking_id", booking.ID),
		zap.Int("user_id", userID),
		zap.Int("cinema_id", req.CinemaID),
		zap.Int("seat_id", req.SeatID))

	return booking, nil
}

// GetUserBookings retrieves all bookings for a user
func (s *BookingService) GetUserBookings(ctx context.Context, userID int) ([]*models.BookingDetail, error) {
	bookings, err := s.bookingRepo.GetUserBookings(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user bookings", zap.Int("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("failed to get bookings")
	}

	return bookings, nil
}

// ProcessPayment processes payment for a booking
func (s *BookingService) ProcessPayment(ctx context.Context, userID int, req *dto.PaymentRequest) (*models.Booking, error) {
	// Get booking
	booking, err := s.bookingRepo.GetByID(ctx, req.BookingID)
	if err != nil {
		s.logger.Error("Booking not found", zap.Int("booking_id", req.BookingID))
		return nil, fmt.Errorf("booking not found")
	}

	// Verify booking belongs to user
	if booking.UserID != userID {
		s.logger.Warn("User attempting to pay for another user's booking",
			zap.Int("user_id", userID),
			zap.Int("booking_id", req.BookingID))
		return nil, fmt.Errorf("unauthorized")
	}

	// Check if already paid
	if booking.PaymentStatus == "paid" {
		s.logger.Warn("Booking already paid", zap.Int("booking_id", req.BookingID))
		return nil, fmt.Errorf("booking is already paid")
	}

	// Validate payment method
	isValid, err := s.paymentRepo.ValidatePaymentMethod(ctx, req.PaymentMethod)
	if err != nil || !isValid {
		s.logger.Error("Invalid payment method", zap.Int("payment_method_id", req.PaymentMethod))
		return nil, fmt.Errorf("invalid payment method")
	}

	// In a real application, you would integrate with a payment gateway here
	// For this example, we'll just update the status to "paid"

	// Update payment status
	if err := s.bookingRepo.UpdatePaymentStatus(ctx, req.BookingID, "paid"); err != nil {
		s.logger.Error("Failed to update payment status", zap.Error(err))
		return nil, fmt.Errorf("failed to process payment")
	}

	// Update booking status
	if err := s.bookingRepo.UpdateBookingStatus(ctx, req.BookingID, "paid"); err != nil {
		s.logger.Error("Failed to update booking status", zap.Error(err))
		return nil, fmt.Errorf("failed to process payment")
	}

	s.logger.Info("Payment processed successfully",
		zap.Int("booking_id", req.BookingID),
		zap.Int("user_id", userID))

	// Get updated booking
	updatedBooking, _ := s.bookingRepo.GetByID(ctx, req.BookingID)
	return updatedBooking, nil
}
