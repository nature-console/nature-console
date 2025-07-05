# Nature Console

A full-stack web application built with Next.js frontend and Go backend using clean architecture.

## System Architecture

- **Frontend**: Next.js with TypeScript and TailwindCSS (Port 3000)
- **Backend**: Go API with clean architecture (Port 8080)
- **Database**: PostgreSQL (Development: Docker, Production: Neon)
- **Deployment**: Fly.io (Production & Staging)

## Development Setup

### Prerequisites

- Node.js 18+
- Go 1.23+
- Docker and Docker Compose

### Quick Start

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd nature-console
   ```

2. **Install dependencies**
   ```bash
   make install
   ```

3. **Start development environment**
   ```bash
   # Start backend and database
   make dev
   
   # In another terminal, start frontend
   make frontend-dev
   ```

4. **Access the application**
   - Frontend: http://localhost:3000
   - API: http://localhost:8080
   - Health check: http://localhost:8080/health
   - Admin login: （メールアドレス・パスワードは .env.development の ADMIN_EMAIL / ADMIN_PASSWORD を参照）

5. **Useful commands**
   ```bash
   make help          # Show all available commands
   make status        # Check container status
   make dev-logs      # View development logs
   make dev-down      # Stop development environment
   ```

## Project Structure

```
nature-console/
├── frontend/                 # Next.js frontend
│   ├── src/
│   │   ├── app/             # App router pages
│   │   ├── components/      # React components
│   │   ├── lib/            # API client and utilities
│   │   └── types/          # TypeScript type definitions
│   ├── package.json
│   └── next.config.js
├── backend/                 # Go backend
│   ├── cmd/                # Application entry point
│   ├── internal/           # Internal application code
│   │   ├── domain/         # Business entities and interfaces
│   │   ├── usecase/        # Business logic
│   │   ├── infrastructure/ # Database and external services
│   │   └── interface/      # HTTP handlers and routing
│   ├── pkg/               # Shared packages
│   ├── Dockerfile
│   └── go.mod
├── docker-compose.yml      # Development environment
├── .env.development       # Development environment variables
├── .env.staging          # Staging environment variables
├── .env.production       # Production environment variables
└── .github/workflows/    # CI/CD configuration
```

## API Endpoints

### Health Check
- `GET /health` - API health status

### Public Articles
- `GET /api/v1/articles` - Get published articles
- `GET /api/v1/articles/:id` - Get article by ID (published only)

### Authentication
- `POST /api/v1/auth/login` - Admin login
- `POST /api/v1/auth/logout` - Admin logout
- `GET /api/v1/auth/me` - Get current admin user (protected)

### Admin Articles (Protected)
- `GET /api/v1/admin/dashboard` - Dashboard statistics
- `GET /api/v1/admin/articles` - Get all articles
- `GET /api/v1/admin/articles/:id` - Get article by ID
- `POST /api/v1/admin/articles` - Create new article
- `PUT /api/v1/admin/articles/:id` - Update article
- `DELETE /api/v1/admin/articles/:id` - Delete article

## Environment Variables

### Development (.env.development)
```env
DATABASE_URL=postgres://postgres:password@db:5432/nature_console_dev?sslmode=disable
API_BASE_URL=http://localhost:8080
PORT=8080
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

### Production (.env.production)
```env
DATABASE_URL=<NEON_PRODUCTION_DB_URL>
API_BASE_URL=https://nature-console.com
PORT=8080
NEXT_PUBLIC_API_BASE_URL=https://nature-console.com
```

## Development Commands

### Development Environment
```bash
# Start development environment (backend + database)
make dev

# Start frontend development server (in another terminal)
make frontend-dev

# Stop development environment
make dev-down

# View logs
make dev-logs

# Check status
make status
```

### Test Environment
```bash
# Start test environment
make test-env

# Stop test environment
make test-down

# View test logs
make test-logs

# Run tests locally
make test
```

### Other Commands
```bash
# Install all dependencies
make install

# Build applications
make build

# Check API health
make health

# Clean all build artifacts and volumes
make clean

# Reset database
make db-reset

# Show all available commands
make help
```

## Database Setup

### Development
The development database runs automatically with Docker Compose.

### Production
1. Create a PostgreSQL database on Neon
2. Update the `DATABASE_URL` in environment variables
3. The application will automatically create required tables

## Deployment

### Staging
Push to `develop` branch to deploy to staging environment.

### Production
Push to `main` branch to deploy to production environment.

### Manual Deployment
```bash
# Install Fly CLI
# Set up Fly.io apps (production and staging)

# Deploy to production
flyctl deploy --config fly.production.toml

# Deploy to staging
flyctl deploy --config fly.staging.toml
```

## Configuration

### Fly.io Setup
1. Install Fly CLI
2. Create apps:
   ```bash
   flyctl apps create nature-console-prod
   flyctl apps create nature-console-staging
   ```
3. Set secrets:
   ```bash
   flyctl secrets set DATABASE_URL=<neon-url> -a nature-console-prod
   flyctl secrets set DATABASE_URL=<neon-url> -a nature-console-staging
   ```

### GitHub Actions
Set the following secrets in your GitHub repository:
- `FLY_API_TOKEN` - Fly.io API token for deployments

## Features

### Backend
- Clean architecture with clear separation of concerns
- JWT-based authentication with HTTP-only cookies
- GORM for database operations
- Database migrations with golang-migrate
- CORS configuration for frontend integration
- Admin user seeding

### Frontend
- Next.js 15 with App Router
- TypeScript for type safety
- TailwindCSS with custom Nature Console design system
- Authentication context with React Context API
- Protected routes with middleware
- Responsive design optimized for blog content

### Full-Stack Integration
- Cookie-based authentication
- Admin dashboard with article statistics
- Complete article management (CRUD operations)
- Public article viewing with published-only filtering
- Landing page showcasing Nature Console mission

### DevOps
- Docker containerization for easy development
- Separate test database configuration
- CI/CD pipeline with GitHub Actions
- Production-ready deployment configuration

## Contributing

1. Create a feature branch from `develop`
2. Make your changes
3. Test thoroughly
4. Submit a pull request to `develop`

## License

This project is licensed under the MIT License.
