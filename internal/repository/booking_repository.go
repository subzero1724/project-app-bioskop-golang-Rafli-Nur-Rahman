package repository

import (
	"context"
	"fmt"

	"cinema-booking-system/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BookingRepository handles booking-related database operations
type BookingRepository struct {
	db *pgxpool.Pool
}

// NewBookingRepository creates a new booking repository
func NewBookingRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{db: db}
}

// Create inserts a new booking into the database
func (r *BookingRepository) Create(ctx context.Context, booking *models.Booking) error {
	query := `
		INSERT INTO bookings (
			user_id, cinema_id, seat_id, booking_date, booking_time,
			payment_method_id, payment_status, total_amount, booking_status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		booking.UserID,
		booking.CinemaID,
		booking.SeatID,
		booking.BookingDate,
		booking.BookingTime,
		booking.PaymentMethodID,
		booking.PaymentStatus,
		booking.TotalAmount,
		booking.BookingStatus,
	).Scan(&booking.ID, &booking.CreatedAt, &booking.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}

	return nil
}

// GetByID retrieves a booking by ID
func (r *BookingRepository) GetByID(ctx context.Context, id int) (*models.Booking, error) {
	query := `
		SELECT id, user_id, cinema_id, seat_id, booking_date, booking_time,
			   payment_method_id, payment_status, total_amount, booking_status,
			   created_at, updated_at
		FROM bookings
		WHERE id = $1
	`

	var booking models.Booking
	err := r.db.QueryRow(ctx, query, id).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.CinemaID,
		&booking.SeatID,
		&booking.BookingDate,
		&booking.BookingTime,
		&booking.PaymentMethodID,
		&booking.PaymentStatus,
		&booking.TotalAmount,
		&booking.BookingStatus,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("booking not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return &booking, nil
}

// CheckSeatAvailability checks if a seat is available for booking
func (r *BookingRepository) CheckSeatAvailability(ctx context.Context, cinemaID, seatID int, date, time string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM bookings
		WHERE cinema_id = $1 
		  AND seat_id = $2 
		  AND booking_date = $3 
		  AND booking_time = $4
		  AND booking_status IN ('reserved', 'paid')
	`

	var count int
	err := r.db.QueryRow(ctx, query, cinemaID, seatID, date, time).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check seat availability: %w", err)
	}

	return count == 0, nil
}

// GetUserBookings retrieves all bookings for a user
func (r *BookingRepository) GetUserBookings(ctx context.Context, userID int) ([]*models.BookingDetail, error) {
	query := `
		SELECT 
			b.id, b.user_id, b.cinema_id, b.seat_id, b.booking_date, b.booking_time,
			b.payment_method_id, b.payment_status, b.total_amount, b.booking_status,
			b.created_at, b.updated_at,
			c.name as cinema_name,
			c.location as cinema_location,
			s.seat_number,
			s.seat_type,
			pm.name as payment_method_name
		FROM bookings b
		INNER JOIN cinemas c ON b.cinema_id = c.id
		INNER JOIN seats s ON b.seat_id = s.id
		LEFT JOIN payment_methods pm ON b.payment_method_id = pm.id
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user bookings: %w", err)
	}
	defer rows.Close()

	var bookings []*models.BookingDetail
	for rows.Next() {
		var booking models.BookingDetail
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.CinemaID,
			&booking.SeatID,
			&booking.BookingDate,
			&booking.BookingTime,
			&booking.PaymentMethodID,
			&booking.PaymentStatus,
			&booking.TotalAmount,
			&booking.BookingStatus,
			&booking.CreatedAt,
			&booking.UpdatedAt,
			&booking.CinemaName,
			&booking.CinemaLocation,
			&booking.SeatNumber,
			&booking.SeatType,
			&booking.PaymentMethodName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, &booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating bookings: %w", err)
	}

	return bookings, nil
}

// UpdatePaymentStatus updates the payment status of a booking
func (r *BookingRepository) UpdatePaymentStatus(ctx context.Context, bookingID int, status string) error {
	query := `
		UPDATE bookings
		SET payment_status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(ctx, query, status, bookingID)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}

// UpdateBookingStatus updates the booking status
func (r *BookingRepository) UpdateBookingStatus(ctx context.Context, bookingID int, status string) error {
	query := `
		UPDATE bookings
		SET booking_status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(ctx, query, status, bookingID)
	if err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}
