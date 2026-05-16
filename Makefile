# Makefile - Common development commands

.PHONY: help build run test clean docker db-setup install-deps

help:
	@echo "Available commands:"
	@echo "  make install-deps  - Install Go dependencies"
	@echo "  make build         - Build application"
	@echo "  make run           - Run application"
	@echo "  make test          - Run tests"
	@echo "  make test-cov      - Run tests with coverage"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run Docker container"
	@echo "  make db-setup      - Setup PostgreSQL database"
	@echo "  make lint          - Run linter"
	@echo "  make format        - Format code"

install-deps:
	@echo "Installing dependencies..."
	go mod download
	go mod verify

build:
	@echo "Building application..."
	go build -o booking-system
	@echo "✓ Build complete: ./booking-system"

run: build
	@echo "Starting application..."
	./booking-system

test:
	@echo "Running tests..."
	go test -v ./...

test-cov:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -f booking-system
	rm -f coverage.out coverage.html
	rm -rf build/
	@echo "✓ Clean complete"

lint:
	@echo "Running linter..."
	golangci-lint run

format:
	@echo "Formatting code..."
	go fmt ./...

docker-build:
	@echo "Building Docker image..."
	docker build -t booking-system:latest .
	@echo "✓ Docker image built: booking-system:latest"

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 booking-system:latest

db-setup:
	@echo "Setting up PostgreSQL database..."
	createdb booking_system || true
	@echo "✓ Database created: booking_system"

db-drop:
	@echo "Dropping PostgreSQL database..."
	dropdb booking_system
	@echo "✓ Database dropped"

logs:
	@echo "Viewing application logs..."
	tail -f logs/app.log

logs-error:
	@echo "Viewing error logs..."
	tail -f logs/error.log

logs-audit:
	@echo "Viewing audit logs..."
	tail -f logs/audit.log

all: clean install-deps test build
	@echo "✓ Full build complete"
