package dto

// RegisterRequest represents user registration input
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"omitempty,max=100"`
}

// LoginRequest represents user login input
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents successful login output
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
	User      UserResponse
}

// UserResponse represents user information in responses
type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name,omitempty"`
}

// BookingRequest represents seat booking input
type BookingRequest struct {
	CinemaID      int    `json:"cinema_id" validate:"required"`
	SeatID        int    `json:"seat_id" validate:"required"`
	Date          string `json:"date" validate:"required"` // YYYY-MM-DD
	Time          string `json:"time" validate:"required"` // HH:MM
	PaymentMethod int    `json:"payment_method" validate:"required"`
}

// PaymentRequest represents payment processing input
type PaymentRequest struct {
	BookingID      int                    `json:"booking_id" validate:"required"`
	PaymentMethod  int                    `json:"payment_method" validate:"required"`
	PaymentDetails map[string]interface{} `json:"payment_details,omitempty"`
}

// PaginationParams represents pagination query parameters
type PaginationParams struct {
	Page     int `validate:"omitempty,min=1"`
	PageSize int `validate:"omitempty,min=1,max=100"`
}

// PaginatedResponse represents paginated API response
type PaginatedResponse struct {
	Data       interface{}      `json:"data"`
	Pagination PaginationMeta   `json:"pagination"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
