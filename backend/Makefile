# Nature Console Backend Makefile

# Default target
.DEFAULT_GOAL := help

# Variables
DOCKER_COMPOSE = docker-compose
GO_CMD = go
SEED_CMD = $(GO_CMD) run cmd/seed/main.go

# Help target
.PHONY: help
help: ## Show this help message
	@echo "Nature Console Backend Commands"
	@echo "=============================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development
.PHONY: dev-up
dev-up: ## Start development environment
	$(DOCKER_COMPOSE) --profile dev up -d

.PHONY: dev-down
dev-down: ## Stop development environment
	$(DOCKER_COMPOSE) --profile dev down

.PHONY: dev-logs
dev-logs: ## Show development logs
	$(DOCKER_COMPOSE) --profile dev logs -f

.PHONY: dev-restart
dev-restart: dev-down dev-up ## Restart development environment

# Testing
.PHONY: test-up
test-up: ## Start test environment
	$(DOCKER_COMPOSE) --profile test up -d

.PHONY: test-down
test-down: ## Stop test environment
	$(DOCKER_COMPOSE) --profile test down

.PHONY: test
test: ## Run tests
	$(GO_CMD) test ./...

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	$(GO_CMD) test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	$(GO_CMD) test -cover ./...

.PHONY: test-race
test-race: ## Run tests with race detection
	$(GO_CMD) test -race ./...

# Building
.PHONY: build
build: ## Build the application
	$(GO_CMD) build -o bin/api cmd/api/main.go

.PHONY: build-seed
build-seed: ## Build the seed command
	$(GO_CMD) build -o bin/seed cmd/seed/main.go

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/

# Database Management
.PHONY: seed
seed: ## Run database seeds (development mode)
	$(SEED_CMD)

.PHONY: seed-prod
seed-prod: ## Run database seeds (production mode)
	$(SEED_CMD) -prod

.PHONY: seed-force
seed-force: ## Force re-seed database
	$(SEED_CMD) -force

.PHONY: seed-help
seed-help: ## Show seed command help
	$(SEED_CMD) -help

# Database seeds with Docker
.PHONY: docker-seed
docker-seed: ## Run seeds inside Docker container
	$(DOCKER_COMPOSE) exec api $(SEED_CMD)

.PHONY: docker-seed-prod
docker-seed-prod: ## Run production seeds inside Docker container
	$(DOCKER_COMPOSE) exec api $(SEED_CMD) -prod

.PHONY: docker-seed-force
docker-seed-force: ## Force re-seed inside Docker container
	$(DOCKER_COMPOSE) exec api $(SEED_CMD) -force

# Code Quality
.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: fmt
fmt: ## Format code
	$(GO_CMD) fmt ./...

.PHONY: vet
vet: ## Run go vet
	$(GO_CMD) vet ./...

.PHONY: mod-tidy
mod-tidy: ## Clean up go.mod
	$(GO_CMD) mod tidy

.PHONY: mod-verify
mod-verify: ## Verify dependencies
	$(GO_CMD) mod verify

# All-in-one commands
.PHONY: check
check: fmt vet test ## Run all checks (format, vet, test)

.PHONY: ci
ci: mod-verify vet test-race test-coverage ## Run CI pipeline

.PHONY: setup
setup: mod-tidy build build-seed ## Setup development environment

# Docker commands
.PHONY: docker-build
docker-build: ## Build Docker image
	$(DOCKER_COMPOSE) build

.PHONY: docker-rebuild
docker-rebuild: ## Rebuild Docker image without cache
	$(DOCKER_COMPOSE) build --no-cache

# Environment commands
.PHONY: env-example
env-example: ## Show example environment variables
	@echo "# Nature Console Environment Variables"
	@echo "DATABASE_URL=postgres://user:password@localhost:5432/nature_console?sslmode=disable"
	@echo "ADMIN_EMAIL=admin@example.com"
	@echo "ADMIN_PASSWORD=your-secure-password"
	@echo "PORT=8080"
	@echo "ENV=development"
	@echo "RUN_SEEDS=true"
	@echo "SEED_MODE=development"

# Development workflow
.PHONY: dev-reset
dev-reset: dev-down clean docker-build dev-up seed ## Complete development reset

.PHONY: quick-start
quick-start: dev-up seed ## Quick start for new developers
	@echo "✅ Development environment is ready!"
	@echo "🌐 API: http://localhost:8080"
	@echo "📊 Health: http://localhost:8080/health"
	@echo "🗄️  Database: PostgreSQL on localhost:5432"