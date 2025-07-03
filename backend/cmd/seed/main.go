package main

import (
	"flag"
	"log"

	"github.com/nature-console/backend/internal/config"
	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/pkg/database"
	"github.com/nature-console/backend/pkg/seed"
	"gorm.io/gorm"
)

func main() {
	// Parse command line flags
	var (
		prodMode = flag.Bool("prod", false, "Run only production seeds (admin users only)")
		force    = flag.Bool("force", false, "Force re-seeding even if data exists")
		help     = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	log.Println("Starting Nature Console Database Seeding...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.New(database.Config{
		URL:        cfg.Database.URL,
		MaxRetries: cfg.Database.MaxRetries,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Auto migrate the schema
	log.Println("Running database migrations...")
	if err := db.AutoMigrate(&entity.Article{}, &entity.AdminUser{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrations completed")

	// Clear existing data if force flag is set
	if *force {
		log.Println("Force flag detected - clearing existing data...")
		if err := clearExistingData(db.DB); err != nil {
			log.Fatalf("Failed to clear existing data: %v", err)
		}
		log.Println("Existing data cleared")
	}

	// Run seeds based on mode
	if *prodMode {
		log.Println("Running in production mode...")
		if err := seed.RunSeedsForProduction(db.DB, cfg.Auth); err != nil {
			log.Fatalf("Failed to run production seeds: %v", err)
		}
	} else {
		log.Println("Running in development mode...")
		if err := seed.RunSeeds(db.DB, cfg.Auth); err != nil {
			log.Fatalf("Failed to run seeds: %v", err)
		}
	}

	log.Println("Database seeding completed successfully!")
}

// clearExistingData removes all existing data when force flag is used
func clearExistingData(db *gorm.DB) error {
	// Delete in reverse order of dependencies
	if err := db.Exec("DELETE FROM articles").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM admin_users").Error; err != nil {
		return err
	}
	return nil
}

// showHelp displays usage information
func showHelp() {
	log.Println(`Nature Console Database Seeder

Usage:
  go run cmd/seed/main.go [flags]

Flags:
  -prod    Run only production seeds (admin users only)
  -force   Force re-seeding even if data exists
  -help    Show this help message

Environment Variables:
  DATABASE_URL     PostgreSQL connection string (required)
  ADMIN_EMAIL      Admin user email (required)
  ADMIN_PASSWORD   Admin user password (required)

Examples:
  # Run development seeds
  go run cmd/seed/main.go

  # Run production seeds only
  go run cmd/seed/main.go -prod

  # Force re-seeding
  go run cmd/seed/main.go -force

  # Run with Docker
  docker-compose exec api go run cmd/seed/main.go`)
}