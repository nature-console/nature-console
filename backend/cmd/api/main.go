package main

import (
	"log"

	"github.com/nature-console/backend/internal/app"
)

func main() {
	// Create and initialize the application
	application, err := app.New()
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}
	defer application.Close()

	// Start the server
	if err := application.Run(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}