package seed

import (
	"errors"
	"log"

	"github.com/nature-console/backend/internal/utils"
	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/internal/config"
	"gorm.io/gorm"
)

// SeedAdminUsers creates initial admin users from configuration
func SeedAdminUsers(db *gorm.DB, authConfig config.AuthConfig) error {
	// Check if admin user already exists
	var count int64
	db.Model(&entity.AdminUser{}).Count(&count)
	if count > 0 {
		log.Println("Admin users already seeded")
		return nil
	}

	// Validate admin credentials
	if authConfig.AdminEmail == "" {
		return errors.New("admin email is not configured")
	}
	if authConfig.AdminPassword == "" {
		return errors.New("admin password is not configured")
	}

	hashedPassword, err := utils.HashPassword(authConfig.AdminPassword)
	if err != nil {
		return err
	}

	adminUser := &entity.AdminUser{
		Email:        authConfig.AdminEmail,
		PasswordHash: hashedPassword,
		Name:         "Nature Console Admin",
	}

	if err := db.Create(adminUser).Error; err != nil {
		return err
	}

	log.Printf("Admin user created: %s", adminUser.Email)
	return nil
}

// SeedArticles creates sample articles for development/testing
func SeedArticles(db *gorm.DB) error {
	// Check if articles already exist
	var count int64
	db.Model(&entity.Article{}).Count(&count)
	if count > 0 {
		log.Println("Articles already seeded")
		return nil
	}

	sampleArticles := []*entity.Article{
		{
			Title:     "Welcome to Nature Console",
			Content:   "This is your first article in Nature Console. You can edit or delete this article from the admin panel.",
			Author:    "Nature Console Admin",
			Published: true,
		},
		{
			Title:     "Getting Started Guide",
			Content:   "Learn how to use Nature Console to manage your blog content effectively.",
			Author:    "Nature Console Admin",
			Published: false,
		},
		{
			Title:     "Nature Photography Tips",
			Content:   "Discover the best techniques for capturing stunning nature photographs.",
			Author:    "Nature Console Admin",
			Published: true,
		},
	}

	for _, article := range sampleArticles {
		if err := db.Create(article).Error; err != nil {
			return err
		}
		log.Printf("Sample article created: %s", article.Title)
	}

	return nil
}

// RunSeeds executes all seeding operations
func RunSeeds(db *gorm.DB, authConfig config.AuthConfig) error {
	log.Println("Running database seeds...")

	if err := SeedAdminUsers(db, authConfig); err != nil {
		return err
	}

	if err := SeedArticles(db); err != nil {
		return err
	}

	log.Println("Database seeding completed")
	return nil
}

// RunSeedsForProduction runs only essential seeds for production
func RunSeedsForProduction(db *gorm.DB, authConfig config.AuthConfig) error {
	log.Println("Running production database seeds...")

	if err := SeedAdminUsers(db, authConfig); err != nil {
		return err
	}

	log.Println("Production database seeding completed")
	return nil
}