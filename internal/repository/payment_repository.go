package repository

import (
	"context"
	"fmt"

	"cinema-booking-system/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PaymentRepository handles payment-related database operations
type PaymentRepository struct {
	db *pgxpool.Pool
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// GetAllPaymentMethods retrieves all active payment methods
func (r *PaymentRepository) GetAllPaymentMethods(ctx context.Context) ([]*models.PaymentMethod, error) {
	query := `
		SELECT id, name, description, is_active, created_at
		FROM payment_methods
		WHERE is_active = true
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}
	defer rows.Close()

	var methods []*models.PaymentMethod
	for rows.Next() {
		var method models.PaymentMethod
		err := rows.Scan(
			&method.ID,
			&method.Name,
			&method.Description,
			&method.IsActive,
			&method.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment method: %w", err)
		}
		methods = append(methods, &method)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating payment methods: %w", err)
	}

	return methods, nil
}

// GetPaymentMethodByID retrieves a payment method by ID
func (r *PaymentRepository) GetPaymentMethodByID(ctx context.Context, id int) (*models.PaymentMethod, error) {
	query := `
		SELECT id, name, description, is_active, created_at
		FROM payment_methods
		WHERE id = $1
	`

	var method models.PaymentMethod
	err := r.db.QueryRow(ctx, query, id).Scan(
		&method.ID,
		&method.Name,
		&method.Description,
		&method.IsActive,
		&method.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("payment method not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}

	return &method, nil
}

// ValidatePaymentMethod checks if a payment method exists and is active
func (r *PaymentRepository) ValidatePaymentMethod(ctx context.Context, id int) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM payment_methods
		WHERE id = $1 AND is_active = true
	`

	var count int
	err := r.db.QueryRow(ctx, query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to validate payment method: %w", err)
	}

	return count > 0, nil
}
