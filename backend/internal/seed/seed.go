package seed

import (
	"errors"
	"log"

	"github.com/nature-console/backend/internal/utils"
	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/internal/config"
	"gorm.io/gorm"
)

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

func RunSeeds(db *gorm.DB, authConfig config.AuthConfig) error {
	log.Println("Running database seeds...")

	if err := SeedAdminUsers(db, authConfig); err != nil {
		return err
	}

	log.Println("Database seeding completed")
	return nil
}
