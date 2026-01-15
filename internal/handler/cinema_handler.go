package handler

import (
	"net/http"
	"strconv"

	"cinema-booking-system/internal/service"
	"cinema-booking-system/internal/utils"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// CinemaHandler handles cinema-related HTTP requests
type CinemaHandler struct {
	cinemaService *service.CinemaService
	logger        *zap.Logger
}

// NewCinemaHandler creates a new cinema handler
func NewCinemaHandler(cinemaService *service.CinemaService, logger *zap.Logger) *CinemaHandler {
	return &CinemaHandler{
		cinemaService: cinemaService,
		logger:        logger,
	}
}

// GetAllCinemas retrieves all cinemas with pagination
// GET /api/cinemas?page=1&page_size=10
func (h *CinemaHandler) GetAllCinemas(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	// Set defaults
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Get cinemas
	result, err := h.cinemaService.GetAllCinemas(r.Context(), page, pageSize)
	if err != nil {
		h.logger.Error("Failed to get cinemas", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get cinemas")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, result)
}

// GetCinemaByID retrieves a specific cinema by ID
// GET /api/cinemas/{cinemaId}
func (h *CinemaHandler) GetCinemaByID(w http.ResponseWriter, r *http.Request) {
	// Get cinema ID from URL
	cinemaIDStr := chi.URLParam(r, "cinemaId")
	cinemaID, err := strconv.Atoi(cinemaIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid cinema ID")
		return
	}

	// Get cinema
	cinema, err := h.cinemaService.GetCinemaByID(r.Context(), cinemaID)
	if err != nil {
		h.logger.Error("Failed to get cinema", zap.Int("cinema_id", cinemaID), zap.Error(err))
		utils.RespondWithError(w, http.StatusNotFound, "Cinema not found")
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, cinema, "")
}

// GetSeatsAvailability retrieves seat availability for a specific showtime
// GET /api/cinemas/{cinemaId}/seats?date={date}&time={time}
func (h *CinemaHandler) GetSeatsAvailability(w http.ResponseWriter, r *http.Request) {
	// Get cinema ID from URL
	cinemaIDStr := chi.URLParam(r, "cinemaId")
	cinemaID, err := strconv.Atoi(cinemaIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid cinema ID")
		return
	}

	// Get query parameters
	date := r.URL.Query().Get("date")
	time := r.URL.Query().Get("time")

	// Validate required parameters
	if date == "" || time == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Date and time parameters are required")
		return
	}

	// Get seat availability
	seats, err := h.cinemaService.GetSeatsAvailability(r.Context(), cinemaID, date, time)
	if err != nil {
		h.logger.Error("Failed to get seat availability",
			zap.Int("cinema_id", cinemaID),
			zap.String("date", date),
			zap.String("time", time),
			zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, seats, "")
}
