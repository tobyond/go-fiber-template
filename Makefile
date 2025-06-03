@.PHONY: help setup db-create db-drop db-reset run dev clean test build assets assets-dev assets-prod

# Default target
help:
	@echo "Available commands:"
	@echo "  setup       - Full setup: create database, install dependencies, and build assets"
	@echo "  db-create   - Create the freshgo database if it doesn't exist"
	@echo "  db-drop     - Drop the freshgo database"
	@echo "  db-reset    - Drop and recreate the database"
	@echo "  assets      - Build frontend assets (development)"
	@echo "  assets-prod - Build frontend assets (production)"
	@echo "  assets-dev  - Build frontend assets and watch for changes"
	@echo "  run         - Run the application"
	@echo "  dev         - Run with auto-reload (requires air) and asset watching"
	@echo "  clean       - Clean build artifacts"
	@echo "  test        - Run tests"
	@echo "  build       - Build the application"

# Full setup
setup: db-create
	@echo "Installing Go dependencies..."
	go mod tidy
	@echo "Creating directories..."
	mkdir -p app/controllers app/models app/services app/middleware config routes tests web/templates web/static/css web/static/js web/assets/css web/assets/js
	@echo "Moving templates to web/templates/ if they exist in templates/"
	@if [ -d "templates" ]; then mv templates/* web/templates/ 2>/dev/null || true; rmdir templates 2>/dev/null || true; fi
	@if [ -d "static" ]; then mv static/* web/static/ 2>/dev/null || true; rmdir static 2>/dev/null || true; fi
	@echo "Installing Node.js dependencies..."
	npm install
	@echo "Building frontend assets..."
	npm run build
	@echo "Creating .env file from example..."
	@if [ ! -f .env ]; then cp .env.example .env 2>/dev/null || true; fi
	@echo "Setup complete! Run 'make dev' to start development."

# Create database if it doesn't exist
db-create:
	@echo "Creating database 'freshgo' if it doesn't exist..."
	@psql -U postgres -h localhost -tc "SELECT 1 FROM pg_database WHERE datname = 'freshgo'" | grep -q 1 || \
	psql -U postgres -h localhost -c "CREATE DATABASE freshgo;"
	@echo "Database 'freshgo' is ready."

# Drop database
db-drop:
	@echo "Dropping database 'freshgo'..."
	@psql -U postgres -h localhost -c "DROP DATABASE IF EXISTS freshgo;"
	@echo "Database 'freshgo' dropped."

# Reset database (drop and recreate)
db-reset: db-drop db-create
	@echo "Database 'freshgo' has been reset."

# Frontend asset commands
assets:
	@echo "Building frontend assets..."
	npm run build

assets-prod:
	@echo "Building frontend assets (production)..."
	npm run build:prod

assets-dev:
	@echo "Starting asset watchers..."
	npm run dev &

# Run the application
run: assets
	@echo "Starting Fresh application..."
	go run *.go

# Development mode with auto-reload and asset watching
dev:
	@echo "Starting development mode..."
	@npm run dev &
	@if command -v air > /dev/null; then \
		echo "Starting Fresh in development mode with auto-reload..."; \
		air; \
	else \
		echo "Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Starting without auto-reload..."; \
		go run *.go; \
	fi."; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	go clean
	rm -f fresh
	rm -rf web/static/css/styles.css*
	rm -rf web/static/js/app.js*
	rm -rf node_modules
	@echo "Clean complete."

# Run tests
test:
	@echo "Running tests..."
	go test ./tests/... -v

# Run tests with verbose output  
test-verbose:
	@echo "Running tests with verbose output..."
	go test ./tests/... -v -count=1

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test ./tests/... -v -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run specific test
test-run:
	@echo "Running specific test (use TEST=TestName)..."
	go test ./tests/... -v -run $(TEST)

# Build the application
build: assets-prod
	@echo "Building Fresh application..."
	go build -o fresh *.go
	@echo "Build complete. Binary: ./fresh"

# Install development dependencies
install-dev:
	@echo "Installing development dependencies..."
	go install github.com/cosmtrek/air@latest
	npm install
	@echo "Development dependencies installed."
