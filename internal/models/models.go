package models

import "time"

// User represents a registered customer
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose password hash in JSON
	FullName     string    `json:"full_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Cinema represents a movie theater
type Cinema struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Description string    `json:"description,omitempty"`
	TotalSeats  int       `json:"total_seats"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Seat represents a seat in a cinema
type Seat struct {
	ID         int     `json:"id"`
	CinemaID   int     `json:"cinema_id"`
	SeatNumber string  `json:"seat_number"`
	RowNumber  string  `json:"row_number"`
	SeatType   string  `json:"seat_type"`
	Price      float64 `json:"price"`
}

// SeatAvailability represents seat status for a specific showtime
type SeatAvailability struct {
	Seat
	IsAvailable bool `json:"is_available"`
}

// PaymentMethod represents available payment options
type PaymentMethod struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// Booking represents a ticket reservation
type Booking struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	CinemaID        int       `json:"cinema_id"`
	SeatID          int       `json:"seat_id"`
	BookingDate     string    `json:"booking_date"` // YYYY-MM-DD format
	BookingTime     string    `json:"booking_time"` // HH:MM:SS format
	PaymentMethodID *int      `json:"payment_method_id,omitempty"`
	PaymentStatus   string    `json:"payment_status"`
	TotalAmount     float64   `json:"total_amount"`
	BookingStatus   string    `json:"booking_status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// BookingDetail extends Booking with related information
type BookingDetail struct {
	Booking
	CinemaName        string  `json:"cinema_name"`
	CinemaLocation    string  `json:"cinema_location"`
	SeatNumber        string  `json:"seat_number"`
	SeatType          string  `json:"seat_type"`
	PaymentMethodName *string `json:"payment_method_name,omitempty"`
}

// Token represents an authentication token
type Token struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
