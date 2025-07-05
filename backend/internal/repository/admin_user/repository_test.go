package admin_user

import (
	"context"
	"testing"

	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/test/mocks"
)

// Test functions
func TestMockAdminUserRepository_Create(t *testing.T) {
	repo := mocks.NewMockAdminUserRepository()
	ctx := context.Background()

	user := &entity.AdminUser{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected ID to be set")
	}

	// Test nil user
	err = repo.Create(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil user")
	}
}

func TestMockAdminUserRepository_GetByID(t *testing.T) {
	repo := mocks.NewMockAdminUserRepository()
	ctx := context.Background()

	user := &entity.AdminUser{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test existing user
	found, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if found.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, found.Email)
	}

	// Test non-existing user
	_, err = repo.GetByID(ctx, 999)
	if err == nil {
		t.Error("Expected error for non-existing user")
	}
}

func TestMockAdminUserRepository_GetByEmail(t *testing.T) {
	repo := mocks.NewMockAdminUserRepository()
	ctx := context.Background()

	user := &entity.AdminUser{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test existing email
	found, err := repo.GetByEmail(ctx, user.Email)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if found.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, found.ID)
	}

	// Test non-existing email
	_, err = repo.GetByEmail(ctx, "nonexistent@example.com")
	if err == nil {
		t.Error("Expected error for non-existing email")
	}
}

func TestMockAdminUserRepository_Update(t *testing.T) {
	repo := mocks.NewMockAdminUserRepository()
	ctx := context.Background()

	user := &entity.AdminUser{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Update user
	user.Name = "Updated User"
	err = repo.Update(ctx, user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify update
	found, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	if found.Name != "Updated User" {
		t.Errorf("Expected name 'Updated User', got %s", found.Name)
	}

	// Test nil user
	err = repo.Update(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil user")
	}

	// Test non-existing user
	nonExistentUser := &entity.AdminUser{ID: 999, Name: "Non-existent"}
	err = repo.Update(ctx, nonExistentUser)
	if err == nil {
		t.Error("Expected error for non-existing user")
	}
}

func TestMockAdminUserRepository_Delete(t *testing.T) {
	repo := mocks.NewMockAdminUserRepository()
	ctx := context.Background()

	user := &entity.AdminUser{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Delete user
	err = repo.Delete(ctx, user.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(ctx, user.ID)
	if err == nil {
		t.Error("Expected error for deleted user")
	}

	// Test non-existing user
	err = repo.Delete(ctx, 999)
	if err == nil {
		t.Error("Expected error for non-existing user")
	}
}

func TestMockAdminUserRepository_List(t *testing.T) {
	repo := mocks.NewMockAdminUserRepository()
	ctx := context.Background()

	// Test empty list
	users, err := repo.List(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(users) != 0 {
		t.Errorf("Expected empty list, got %d users", len(users))
	}

	// Add users
	user1 := &entity.AdminUser{Email: "user1@example.com", Name: "User 1"}
	user2 := &entity.AdminUser{Email: "user2@example.com", Name: "User 2"}

	repo.Create(ctx, user1)
	repo.Create(ctx, user2)

	// Test list with users
	users, err = repo.List(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}