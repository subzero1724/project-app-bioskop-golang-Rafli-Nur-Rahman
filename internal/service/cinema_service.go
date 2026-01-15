package service

import (
	"context"
	"fmt"

	"cinema-booking-system/internal/dto"
	"cinema-booking-system/internal/models"
	"cinema-booking-system/internal/repository"

	"go.uber.org/zap"
)

// CinemaService handles cinema-related business logic
type CinemaService struct {
	cinemaRepo *repository.CinemaRepository
	logger     *zap.Logger
}

// NewCinemaService creates a new cinema service
func NewCinemaService(cinemaRepo *repository.CinemaRepository, logger *zap.Logger) *CinemaService {
	return &CinemaService{
		cinemaRepo: cinemaRepo,
		logger:     logger,
	}
}

// GetAllCinemas retrieves all cinemas with pagination
func (s *CinemaService) GetAllCinemas(ctx context.Context, page, pageSize int) (*dto.PaginatedResponse, error) {
	// Calculate offset
	offset := (page - 1) * pageSize

	// Get cinemas
	cinemas, err := s.cinemaRepo.GetAll(ctx, pageSize, offset)
	if err != nil {
		s.logger.Error("Failed to get cinemas", zap.Error(err))
		return nil, fmt.Errorf("failed to get cinemas")
	}

	// Get total count
	totalCount, err := s.cinemaRepo.Count(ctx)
	if err != nil {
		s.logger.Error("Failed to count cinemas", zap.Error(err))
		return nil, fmt.Errorf("failed to count cinemas")
	}

	// Calculate total pages
	totalPages := (totalCount + pageSize - 1) / pageSize

	return &dto.PaginatedResponse{
		Data: cinemas,
		Pagination: dto.PaginationMeta{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalItems:  totalCount,
			TotalPages:  totalPages,
		},
	}, nil
}

// GetCinemaByID retrieves a cinema by ID
func (s *CinemaService) GetCinemaByID(ctx context.Context, id int) (*models.Cinema, error) {
	cinema, err := s.cinemaRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get cinema", zap.Int("cinema_id", id), zap.Error(err))
		return nil, fmt.Errorf("cinema not found")
	}

	return cinema, nil
}

// GetSeatsAvailability retrieves seat availability for a specific showtime
func (s *CinemaService) GetSeatsAvailability(ctx context.Context, cinemaID int, date, time string) ([]*models.SeatAvailability, error) {
	// Validate cinema exists
	_, err := s.cinemaRepo.GetByID(ctx, cinemaID)
	if err != nil {
		s.logger.Error("Cinema not found", zap.Int("cinema_id", cinemaID))
		return nil, fmt.Errorf("cinema not found")
	}

	// Get seat availability
	seats, err := s.cinemaRepo.GetSeatsAvailability(ctx, cinemaID, date, time)
	if err != nil {
		s.logger.Error("Failed to get seat availability",
			zap.Int("cinema_id", cinemaID),
			zap.String("date", date),
			zap.String("time", time),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get seat availability")
	}

	return seats, nil
}
