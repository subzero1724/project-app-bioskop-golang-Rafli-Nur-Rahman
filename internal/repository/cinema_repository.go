package repository

import (
	"context"
	"fmt"

	"cinema-booking-system/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CinemaRepository handles cinema-related database operations
type CinemaRepository struct {
	db *pgxpool.Pool
}

// NewCinemaRepository creates a new cinema repository
func NewCinemaRepository(db *pgxpool.Pool) *CinemaRepository {
	return &CinemaRepository{db: db}
}

// GetAll retrieves all cinemas with pagination
func (r *CinemaRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Cinema, error) {
	query := `
		SELECT id, name, location, description, total_seats, created_at, updated_at
		FROM cinemas
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get cinemas: %w", err)
	}
	defer rows.Close()

	var cinemas []*models.Cinema
	for rows.Next() {
		var cinema models.Cinema
		err := rows.Scan(
			&cinema.ID,
			&cinema.Name,
			&cinema.Location,
			&cinema.Description,
			&cinema.TotalSeats,
			&cinema.CreatedAt,
			&cinema.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cinema: %w", err)
		}
		cinemas = append(cinemas, &cinema)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cinemas: %w", err)
	}

	return cinemas, nil
}

// Count returns the total number of cinemas
func (r *CinemaRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM cinemas`

	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count cinemas: %w", err)
	}

	return count, nil
}

// GetByID retrieves a cinema by ID
func (r *CinemaRepository) GetByID(ctx context.Context, id int) (*models.Cinema, error) {
	query := `
		SELECT id, name, location, description, total_seats, created_at, updated_at
		FROM cinemas
		WHERE id = $1
	`

	var cinema models.Cinema
	err := r.db.QueryRow(ctx, query, id).Scan(
		&cinema.ID,
		&cinema.Name,
		&cinema.Location,
		&cinema.Description,
		&cinema.TotalSeats,
		&cinema.CreatedAt,
		&cinema.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("cinema not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cinema: %w", err)
	}

	return &cinema, nil
}

// GetSeats retrieves all seats for a cinema
func (r *CinemaRepository) GetSeats(ctx context.Context, cinemaID int) ([]*models.Seat, error) {
	query := `
		SELECT id, cinema_id, seat_number, row_number, seat_type, price
		FROM seats
		WHERE cinema_id = $1
		ORDER BY row_number, seat_number
	`

	rows, err := r.db.Query(ctx, query, cinemaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seats: %w", err)
	}
	defer rows.Close()

	var seats []*models.Seat
	for rows.Next() {
		var seat models.Seat
		err := rows.Scan(
			&seat.ID,
			&seat.CinemaID,
			&seat.SeatNumber,
			&seat.RowNumber,
			&seat.SeatType,
			&seat.Price,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan seat: %w", err)
		}
		seats = append(seats, &seat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating seats: %w", err)
	}

	return seats, nil
}

// GetSeatsAvailability retrieves seat availability for a specific showtime
func (r *CinemaRepository) GetSeatsAvailability(ctx context.Context, cinemaID int, date, time string) ([]*models.SeatAvailability, error) {
	query := `
		SELECT 
			s.id,
			s.cinema_id,
			s.seat_number,
			s.row_number,
			s.seat_type,
			s.price,
			CASE 
				WHEN b.id IS NULL THEN true 
				ELSE false 
			END as is_available
		FROM seats s
		LEFT JOIN bookings b ON s.id = b.seat_id 
			AND b.cinema_id = $1 
			AND b.booking_date = $2 
			AND b.booking_time = $3
			AND b.booking_status IN ('reserved', 'paid')
		WHERE s.cinema_id = $1
		ORDER BY s.row_number, s.seat_number
	`

	rows, err := r.db.Query(ctx, query, cinemaID, date, time)
	if err != nil {
		return nil, fmt.Errorf("failed to get seat availability: %w", err)
	}
	defer rows.Close()

	var seats []*models.SeatAvailability
	for rows.Next() {
		var seat models.SeatAvailability
		err := rows.Scan(
			&seat.ID,
			&seat.CinemaID,
			&seat.SeatNumber,
			&seat.RowNumber,
			&seat.SeatType,
			&seat.Price,
			&seat.IsAvailable,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan seat availability: %w", err)
		}
		seats = append(seats, &seat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating seat availability: %w", err)
	}

	return seats, nil
}

// GetSeatByID retrieves a seat by ID
func (r *CinemaRepository) GetSeatByID(ctx context.Context, seatID int) (*models.Seat, error) {
	query := `
		SELECT id, cinema_id, seat_number, row_number, seat_type, price
		FROM seats
		WHERE id = $1
	`

	var seat models.Seat
	err := r.db.QueryRow(ctx, query, seatID).Scan(
		&seat.ID,
		&seat.CinemaID,
		&seat.SeatNumber,
		&seat.RowNumber,
		&seat.SeatType,
		&seat.Price,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("seat not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get seat: %w", err)
	}

	return &seat, nil
}
