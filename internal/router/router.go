package router

import (
	"cinema-booking-system/internal/handler"
	"cinema-booking-system/internal/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

// SetupRouter configures all API routes
func SetupRouter(
	authHandler *handler.AuthHandler,
	cinemaHandler *handler.CinemaHandler,
	bookingHandler *handler.BookingHandler,
	paymentHandler *handler.PaymentHandler,
	authMiddleware *middleware.AuthMiddleware,
	loggingMiddleware *middleware.LoggingMiddleware,
) *chi.Mux {
	r := chi.NewRouter()

	// Global middlewares
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Recoverer)
	r.Use(loggingMiddleware.Log)

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Public routes (no authentication required)
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Get("/cinemas", cinemaHandler.GetAllCinemas)
		r.Get("/cinemas/{cinemaId}", cinemaHandler.GetCinemaByID)
		r.Get("/cinemas/{cinemaId}/seats", cinemaHandler.GetSeatsAvailability)
		r.Get("/payment-methods", paymentHandler.GetAllPaymentMethods)

		// Protected routes (authentication required)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			// Authentication
			r.Post("/logout", authHandler.Logout)

			// Booking
			r.Post("/booking", bookingHandler.CreateBooking)
			r.Get("/user/bookings", bookingHandler.GetUserBookings)

			// Payment
			r.Post("/pay", bookingHandler.ProcessPayment)
		})
	})

	return r
}
