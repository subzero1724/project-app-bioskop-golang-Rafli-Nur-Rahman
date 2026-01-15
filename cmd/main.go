package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cinema-booking-system/internal/config"
	"cinema-booking-system/internal/database"
	"cinema-booking-system/internal/handler"
	"cinema-booking-system/internal/middleware"
	"cinema-booking-system/internal/repository"
	"cinema-booking-system/internal/router"
	"cinema-booking-system/internal/service"
	"cinema-booking-system/internal/utils"
	"cinema-booking-system/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.NewLogger(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting Cinema Booking System",
		zap.String("app_name", cfg.App.Name),
		zap.String("environment", cfg.App.Env),
		zap.String("port", cfg.App.Port),
	)

	// Connect to database
	db, err := database.NewConnection(cfg, log)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	cinemaRepo := repository.NewCinemaRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg, log)
	cinemaService := service.NewCinemaService(cinemaRepo, log)
	bookingService := service.NewBookingService(bookingRepo, cinemaRepo, paymentRepo, log)
	paymentService := service.NewPaymentService(paymentRepo, log)

	// Initialize validator
	validator := utils.NewValidator()

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, validator, log)
	cinemaHandler := handler.NewCinemaHandler(cinemaService, log)
	bookingHandler := handler.NewBookingHandler(bookingService, validator, log)
	paymentHandler := handler.NewPaymentHandler(paymentService, log)

	// Initialize middlewares
	authMiddleware := middleware.NewAuthMiddleware(authService, log)
	loggingMiddleware := middleware.NewLoggingMiddleware(log)

	// Setup router
	r := router.SetupRouter(
		authHandler,
		cinemaHandler,
		bookingHandler,
		paymentHandler,
		authMiddleware,
		loggingMiddleware,
	)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info("Server starting", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Server is shutting down...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited gracefully")
}
