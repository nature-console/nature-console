package database

import (
	"fmt"
	"log"
	"time"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds database configuration
type Config struct {
	URL        string
	MaxRetries int
}

// Database wraps the GORM database instance
type Database struct {
	*gorm.DB
	config Config
}

// New creates a new database instance
func New(config Config) (*Database, error) {
	if config.MaxRetries == 0 {
		config.MaxRetries = 10
	}

	db, err := connectWithRetry(config.URL, config.MaxRetries)
	if err != nil {
		return nil, err
	}

	return &Database{
		DB:     db,
		config: config,
	}, nil
}

// connectWithRetry attempts to connect to the database with retry logic
func connectWithRetry(dbURL string, maxRetries int) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to connect to database (attempt %d/%d)", i+1, maxRetries)
		
		db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
		if err == nil {
			// Test the connection
			sqlDB, sqlErr := db.DB()
			if sqlErr != nil {
				err = sqlErr
			} else if pingErr := sqlDB.Ping(); pingErr != nil {
				err = pingErr
			} else {
				log.Println("Database connected successfully")
				return db, nil
			}
		}

		if i < maxRetries-1 {
			log.Printf("Failed to connect to database: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
		}
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Health checks if the database connection is healthy
func (d *Database) Health() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}