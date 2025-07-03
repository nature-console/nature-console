.PHONY: help dev test-env build clean test lint docker-up docker-down install dev-down test-down dev-logs test-logs

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install dependencies
	cd frontend && npm install
	cd backend && go mod tidy

# Development environment
dev: ## Start development environment (API + DB only)
	@echo "Starting development environment..."
	docker-compose up api db
	@echo "Backend started on http://localhost:8080"
	@echo "Database available on port 5432"
	@echo ""
	@echo "To start frontend: cd frontend && npm run dev"

dev-down: ## Stop development environment
	docker-compose stop api db
	@echo "Development environment stopped"

dev-logs: ## Show development environment logs
	docker-compose logs -f api db

# Test environment
test-env: ## Start test environment (Test DB + Test runner)
	@echo "Starting test environment..."
	docker-compose --profile test up -d test-db test
	@echo "Test database started on port 5433"
	@echo "Tests will run automatically"

test-down: ## Stop test environment
	docker-compose --profile test stop test-db test
	@echo "Test environment stopped"

test-logs: ## Show test environment logs
	docker-compose --profile test logs -f test-db test

# General commands
docker-up: ## Start all containers (dev + test)
	docker-compose --profile test up -d

docker-down: ## Stop all containers
	docker-compose --profile test down

build: ## Build applications
	cd frontend && npm run build
	cd backend && go build -o main cmd/api/main.go

test: ## Run tests locally
	cd backend && go test ./...

lint: ## Run linters
	cd frontend && npm run lint
	cd backend && go fmt ./...

clean: ## Clean build artifacts and volumes
	cd frontend && rm -rf .next node_modules
	cd backend && rm -f main
	docker-compose --profile test down -v
	docker system prune -f

health: ## Check API health
	@echo "Checking API health..."
	@curl -f http://localhost:8080/health || echo "API is not running. Start with 'make dev'"

status: ## Show container status
	docker-compose ps

# Frontend specific
frontend-dev: ## Start frontend development server
	cd frontend && npm run dev

frontend-build: ## Build frontend for production
	cd frontend && npm run build

# Backend specific
backend-dev: ## Start backend with local database
	cd backend && go run cmd/api/main.go

backend-test: ## Run backend tests
	cd backend && go test ./...

# Database operations
db-reset: ## Reset development database
	docker-compose stop api
	docker-compose rm -f db
	docker volume rm nature-console_postgres_data
	make dev

# Show useful commands
quick-start: ## Show quick start commands
	@echo "=== Nature Console Quick Start ==="
	@echo ""
	@echo "1. Start development environment:"
	@echo "   make dev"
	@echo ""
	@echo "2. Start frontend (in another terminal):"
	@echo "   make frontend-dev"
	@echo ""
	@echo "3. Access application:"
	@echo "   Frontend: http://localhost:3000"
	@echo "   API: http://localhost:8080"
	@echo ""
	@echo "4. Check status:"
	@echo "   make status"
	@echo ""
	@echo "5. View logs:"
	@echo "   make dev-logs"
	@echo ""
	@echo "6. Stop environment:"
	@echo "   make dev-down"