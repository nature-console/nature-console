package app

import (
	"log"

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
	"github.com/nature-console/backend/internal/seed"
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

	// Run seeds
	if err := seed.RunSeeds(a.database.DB, a.config.Auth); err != nil {
		return err
	}

	log.Println("Database connected, migrated, and seeded successfully")

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

// Close gracefully shuts down the application
func (a *App) Close() error {
	return a.database.Close()
}