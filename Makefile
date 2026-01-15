.PHONY: help run build clean test migrate-up migrate-down docker-up docker-down

# Default target
help:
	@echo "Available commands:"
	@echo "  make run          - Run the application"
	@echo "  make build        - Build the application"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"

# Run the application
run:
	@echo "Starting Cinema Booking System..."
	@go run cmd/main.go

# Build the application
build:
	@echo "Building Cinema Booking System..."
	@go build -o bin/cinema-booking-system cmd/main.go
	@echo "Build complete! Binary: bin/cinema-booking-system"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run database migrations
migrate-up:
	@echo "Running database migrations..."
	@psql -U postgres -d cinema_booking -f migrations/001_init_schema.sql
	@echo "Migrations complete!"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@echo "Dependencies installed!"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted!"

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run
	@echo "Linting complete!"

# Create .env file from example
setup-env:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env file created from .env.example"; \
		echo "Please update the values in .env file"; \
	else \
		echo ".env file already exists"; \
	fi
