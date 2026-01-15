package repository

import (
	"context"
	"fmt"
	"time"

	"cinema-booking-system/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository handles user-related database operations
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, full_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// CreateToken stores a new authentication token
func (r *UserRepository) CreateToken(ctx context.Context, token *models.Token) error {
	query := `
		INSERT INTO tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		token.UserID,
		token.Token,
		token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}

	return nil
}

// GetTokenByValue retrieves a token by its value
func (r *UserRepository) GetTokenByValue(ctx context.Context, tokenValue string) (*models.Token, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM tokens
		WHERE token = $1 AND expires_at > $2
	`

	var token models.Token
	err := r.db.QueryRow(ctx, query, tokenValue, time.Now()).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("token not found or expired")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	return &token, nil
}

// DeleteToken removes a token from the database
func (r *UserRepository) DeleteToken(ctx context.Context, tokenValue string) error {
	query := `DELETE FROM tokens WHERE token = $1`

	_, err := r.db.Exec(ctx, query, tokenValue)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}

// DeleteExpiredTokens removes all expired tokens
func (r *UserRepository) DeleteExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM tokens WHERE expires_at < $1`

	_, err := r.db.Exec(ctx, query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete expired tokens: %w", err)
	}

	return nil
}
