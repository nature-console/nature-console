package auth

import (
	"context"
	"testing"

	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/internal/utils"
	"github.com/nature-console/backend/test/mocks"
)

func TestUseCase_Login(t *testing.T) {
	mockRepo := mocks.NewMockAdminUserRepository()
	useCase := NewUseCase(mockRepo)
	ctx := context.Background()

	// Create test user
	password := "testpassword"
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	testUser := &entity.AdminUser{
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		Name:         "Test User",
	}

	err = mockRepo.Create(ctx, testUser)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name        string
		email       string
		password    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid credentials",
			email:       "test@example.com",
			password:    "testpassword",
			expectError: false,
		},
		{
			name:        "Invalid email",
			email:       "invalid@example.com",
			password:    "testpassword",
			expectError: true,
			errorMsg:    "invalid credentials",
		},
		{
			name:        "Invalid password",
			email:       "test@example.com",
			password:    "wrongpassword",
			expectError: true,
			errorMsg:    "invalid credentials",
		},
		{
			name:        "Empty email",
			email:       "",
			password:    "testpassword",
			expectError: true,
			errorMsg:    "email and password are required",
		},
		{
			name:        "Empty password",
			email:       "test@example.com",
			password:    "",
			expectError: true,
			errorMsg:    "email and password are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := useCase.Login(ctx, tt.email, tt.password)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
					return
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
				if response != nil {
					t.Errorf("Expected nil response on error, got %v", response)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
					return
				}
				if response == nil {
					t.Error("Expected response, got nil")
					return
				}
				if response.Token == "" {
					t.Error("Expected token in response")
				}
				if response.User == nil {
					t.Error("Expected user in response")
				} else {
					if response.User.Email != tt.email {
						t.Errorf("Expected user email %s, got %s", tt.email, response.User.Email)
					}
					// Password hash should not be returned
					if response.User.PasswordHash != "" {
						t.Error("Password hash should not be returned in response")
					}
				}
			}
		})
	}
}

func TestUseCase_GetUserFromToken(t *testing.T) {
	mockRepo := mocks.NewMockAdminUserRepository()
	useCase := NewUseCase(mockRepo)
	ctx := context.Background()

	// Create test user
	testUser := &entity.AdminUser{
		Email: "test@example.com",
		Name:  "Test User",
	}

	err := mockRepo.Create(ctx, testUser)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Generate valid token
	validToken, err := utils.GenerateToken(testUser.ID, testUser.Email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	tests := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "Valid token",
			token:       validToken,
			expectError: false,
		},
		{
			name:        "Invalid token",
			token:       "invalid.token.here",
			expectError: true,
		},
		{
			name:        "Empty token",
			token:       "",
			expectError: true,
		},
		{
			name:        "Malformed token",
			token:       "malformed-token",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := useCase.GetUserFromToken(ctx, tt.token)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if user != nil {
					t.Errorf("Expected nil user on error, got %v", user)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
					return
				}
				if user == nil {
					t.Error("Expected user, got nil")
					return
				}
				if user.Email != testUser.Email {
					t.Errorf("Expected user email %s, got %s", testUser.Email, user.Email)
				}
				if user.ID != testUser.ID {
					t.Errorf("Expected user ID %d, got %d", testUser.ID, user.ID)
				}
			}
		})
	}
}

func TestUseCase_GetUserFromToken_UserNotFound(t *testing.T) {
	mockRepo := mocks.NewMockAdminUserRepository()
	useCase := NewUseCase(mockRepo)
	ctx := context.Background()

	// Generate token for non-existent user
	nonExistentUserID := uint(999)
	token, err := utils.GenerateToken(nonExistentUserID, "nonexistent@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	user, err := useCase.GetUserFromToken(ctx, token)
	if err == nil {
		t.Error("Expected error for non-existent user")
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}
}

func TestUseCase_LoginPasswordValidation(t *testing.T) {
	mockRepo := mocks.NewMockAdminUserRepository()
	useCase := NewUseCase(mockRepo)
	ctx := context.Background()

	// Test different password scenarios
	passwords := []string{"short", "medium_length", "very_long_password_with_numbers123"}

	for _, password := range passwords {
		t.Run("Password_"+password, func(t *testing.T) {
			hashedPassword, err := utils.HashPassword(password)
			if err != nil {
				t.Fatalf("Failed to hash password: %v", err)
			}

			testUser := &entity.AdminUser{
				Email:        "test_" + password + "@example.com",
				PasswordHash: hashedPassword,
				Name:         "Test User",
			}

			err = mockRepo.Create(ctx, testUser)
			if err != nil {
				t.Fatalf("Failed to create test user: %v", err)
			}

			// Test correct password
			response, err := useCase.Login(ctx, testUser.Email, password)
			if err != nil {
				t.Errorf("Expected no error for correct password, got %v", err)
			}
			if response == nil {
				t.Error("Expected response for correct password")
			}

			// Test incorrect password
			_, err = useCase.Login(ctx, testUser.Email, password+"wrong")
			if err == nil {
				t.Error("Expected error for incorrect password")
			}
		})
	}
}

func TestUseCase_LoginTokenGeneration(t *testing.T) {
	mockRepo := mocks.NewMockAdminUserRepository()
	useCase := NewUseCase(mockRepo)
	ctx := context.Background()

	password := "testpassword"
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	testUser := &entity.AdminUser{
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		Name:         "Test User",
	}

	err = mockRepo.Create(ctx, testUser)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Login multiple times and ensure different tokens are generated
	response1, err := useCase.Login(ctx, testUser.Email, password)
	if err != nil {
		t.Fatalf("Failed to login first time: %v", err)
	}

	response2, err := useCase.Login(ctx, testUser.Email, password)
	if err != nil {
		t.Fatalf("Failed to login second time: %v", err)
	}

	if response1.Token == response2.Token {
		t.Error("Expected different tokens for multiple logins")
	}

	// Verify tokens are valid
	_, err = utils.ValidateToken(response1.Token)
	if err != nil {
		t.Errorf("First token is invalid: %v", err)
	}

	_, err = utils.ValidateToken(response2.Token)
	if err != nil {
		t.Errorf("Second token is invalid: %v", err)
	}
}