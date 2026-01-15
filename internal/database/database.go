package database

import (
	"context"
	"fmt"

	"cinema-booking-system/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// NewConnection creates a new PostgreSQL connection pool
func NewConnection(cfg *config.Config, logger *zap.Logger) (*pgxpool.Pool, error) {
	ctx := context.Background()

	// Parse connection configuration
	poolConfig, err := pgxpool.ParseConfig(cfg.GetDatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Set pool settings
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established successfully",
		zap.String("host", cfg.Database.Host),
		zap.String("database", cfg.Database.Name),
	)

	return pool, nil
}
