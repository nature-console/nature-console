package app

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nature-console/backend/internal/config"
	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/internal/handler/admin"
	"github.com/nature-console/backend/internal/handler/article"
	"github.com/nature-console/backend/internal/handler/auth"
	"github.com/nature-console/backend/internal/middleware"
	adminUserRepo "github.com/nature-console/backend/internal/repository/admin_user"
	articleRepo "github.com/nature-console/backend/internal/repository/article"
	"github.com/nature-console/backend/internal/routes"
	"github.com/nature-console/backend/pkg/seed"
	articleUC "github.com/nature-console/backend/internal/usecase/article"
	authUC "github.com/nature-console/backend/internal/usecase/auth"
	"github.com/nature-console/backend/pkg/database"
)

// App represents the application
type App struct {
	config   *config.Config
	database *database.Database
	router   *gin.Engine
}

// New creates a new application instance
func New() (*App, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Connect to database
	db, err := database.New(database.Config{
		URL:        cfg.Database.URL,
		MaxRetries: cfg.Database.MaxRetries,
	})
	if err != nil {
		return nil, err
	}

	// Create app instance
	app := &App{
		config:   cfg,
		database: db,
		router:   gin.Default(),
	}

	// Initialize the application
	if err := app.initialize(); err != nil {
		return nil, err
	}

	return app, nil
}

// initialize sets up the application
func (a *App) initialize() error {
	// Auto migrate the schema
	if err := a.database.AutoMigrate(&entity.Article{}, &entity.AdminUser{}); err != nil {
		return err
	}

	// Run seeds only if enabled via environment variable
	if a.shouldRunSeeds() {
		if err := a.runSeeds(); err != nil {
			return err
		}
		log.Println("Database connected, migrated, and seeded successfully")
	} else {
		log.Println("Database connected and migrated successfully")
	}

	// Setup dependencies
	a.setupDependencies()

	// Setup middleware and routes
	a.setupRouter()

	return nil
}

// setupDependencies initializes all application dependencies
func (a *App) setupDependencies() {
	// Initialize repositories
	articleRepository := articleRepo.NewArticleRepository(a.database.DB)
	adminUserRepository := adminUserRepo.NewAdminUserRepository(a.database.DB)

	// Initialize use cases
	articleUseCase := articleUC.NewUseCase(articleRepository)
	authUseCase := authUC.NewUseCase(adminUserRepository)

	// Initialize handlers
	articleHandler := article.NewHandler(articleUseCase)
	authHandler := auth.NewHandler(authUseCase)
	adminHandler := admin.NewHandler(articleUseCase)

	// Store handlers for route setup
	a.setupRoutes(articleHandler, authHandler, adminHandler, authUseCase)
}

// setupRouter configures middleware and routes
func (a *App) setupRouter() {
	// Apply global middleware
	middleware.SetupMiddlewares(a.router)

	// Health check endpoint
	a.router.GET("/health", a.healthCheck)
}

// setupRoutes configures all application routes
func (a *App) setupRoutes(
	articleHandler *article.Handler,
	authHandler *auth.Handler,
	adminHandler *admin.Handler,
	authUseCase *authUC.UseCase,
) {
	// Setup API routes
	api := a.router.Group("/api/v1")

	// Public routes
	routes.SetupArticleRoutes(api, articleHandler)

	// Auth routes
	routes.SetupAuthRoutes(api, authHandler, authUseCase)

	// Admin routes (protected)
	routes.SetupAdminRoutes(api, adminHandler, authUseCase)
}

// healthCheck handles health check requests
func (a *App) healthCheck(c *gin.Context) {
	// Check database health
	if err := a.database.Health(); err != nil {
		c.JSON(500, gin.H{
			"status":  "unhealthy",
			"message": "Database connection failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "healthy",
		"message": "Nature Console API is running",
	})
}

// Run starts the HTTP server
func (a *App) Run() error {
	log.Printf("Server starting on port %s", a.config.Server.Port)
	return a.router.Run(":" + a.config.Server.Port)
}

// shouldRunSeeds determines if seeds should be run based on environment variables
func (a *App) shouldRunSeeds() bool {
	// Check various environment variables that indicate seeding should occur
	runSeeds := os.Getenv("RUN_SEEDS")
	autoSeed := os.Getenv("AUTO_SEED")
	env := os.Getenv("ENV")
	ginMode := os.Getenv("GIN_MODE")

	// Run seeds if explicitly requested
	if strings.ToLower(runSeeds) == "true" || strings.ToLower(autoSeed) == "true" {
		return true
	}

	// Run seeds in development environment by default
	if env == "development" || env == "dev" || ginMode == "debug" {
		// Check if explicitly disabled
		if strings.ToLower(runSeeds) == "false" || strings.ToLower(autoSeed) == "false" {
			return false
		}
		return true
	}

	// Don't run seeds in production by default
	return false
}

// runSeeds executes the appropriate seeding strategy based on environment
func (a *App) runSeeds() error {
	env := os.Getenv("ENV")
	seedMode := os.Getenv("SEED_MODE")

	// Determine seeding mode
	if env == "production" || env == "prod" || seedMode == "production" {
		log.Println("Running production seeds (admin users only)...")
		return seed.RunSeedsForProduction(a.database.DB, a.config.Auth)
	}

	log.Println("Running development seeds (admin users + sample data)...")
	return seed.RunSeeds(a.database.DB, a.config.Auth)
}

// Close gracefully shuts down the application
func (a *App) Close() error {
	return a.database.Close()
}