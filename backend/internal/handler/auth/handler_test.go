package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/internal/domain/repository"
	authUC "github.com/nature-console/backend/internal/usecase/auth"
	"github.com/nature-console/backend/internal/utils"
	"github.com/nature-console/backend/test/mocks"
)

func setupTestHandler() (*Handler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockRepo := mocks.NewMockAdminUserRepository()
	useCase := authUC.NewUseCase(mockRepo)
	handler := NewHandler(useCase)
	
	router := gin.New()
	return handler, router
}

func createTestUser(repo repository.AdminUserRepository, email, password string) (*entity.AdminUser, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &entity.AdminUser{
		Email:        email,
		PasswordHash: hashedPassword,
		Name:         "Test User",
	}

	err = repo.Create(context.Background(), user)
	return user, err
}

func TestHandler_Login(t *testing.T) {
	// Setup shared test data
	mockRepo := mocks.NewMockAdminUserRepository()
	_, err := createTestUser(mockRepo, "test@example.com", "testpassword")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "Valid login",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "testpassword",
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name: "Invalid email",
			requestBody: map[string]string{
				"email":    "invalid@example.com",
				"password": "testpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
		{
			name: "Invalid password",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
		{
			name: "Missing email",
			requestBody: map[string]string{
				"password": "testpassword",
			},
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
		},
		{
			name: "Missing password",
			requestBody: map[string]string{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
		},
		{
			name:           "Empty request body",
			requestBody:    map[string]string{},
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh handler and router for each test
			gin.SetMode(gin.TestMode)
			useCase := authUC.NewUseCase(mockRepo)
			handler := NewHandler(useCase)
			router := gin.New()
			router.POST("/login", handler.Login)

			// Marshal request body
			requestBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			
			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response body
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if tt.expectToken {
				if _, exists := response["user"]; !exists {
					t.Error("Expected user in response")
				}
				if _, exists := response["message"]; !exists {
					t.Error("Expected message in response")
				}
			} else {
				if _, exists := response["error"]; !exists {
					t.Error("Expected error in response")
				}
			}

			// Check cookie for successful login
			if tt.expectedStatus == http.StatusOK && tt.expectToken {
				cookies := w.Result().Cookies()
				found := false
				for _, cookie := range cookies {
					if cookie.Name == "token" {
						found = true
						if cookie.Value == "" {
							t.Error("Token cookie should not be empty")
						}
						if !cookie.HttpOnly {
							t.Error("Token cookie should be HttpOnly")
						}
						break
					}
				}
				if !found {
					t.Error("Expected token cookie to be set")
				}
			}
		})
	}
}

func TestHandler_Logout(t *testing.T) {
	handler, router := setupTestHandler()
	
	// Setup route
	router.POST("/logout", handler.Logout)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if message, exists := response["message"]; !exists || message != "Logout successful" {
		t.Error("Expected logout success message")
	}

	// Check that token cookie is cleared
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			found = true
			if cookie.Value != "" {
				t.Error("Token cookie should be empty after logout")
			}
			if cookie.MaxAge != -1 {
				t.Error("Token cookie should have MaxAge -1 to delete it")
			}
			break
		}
	}
	if !found {
		t.Error("Expected token cookie to be set for deletion")
	}
}

func TestHandler_Me(t *testing.T) {
	handler, router := setupTestHandler()
	
	// Setup route with auth middleware mock
	router.GET("/me", func(c *gin.Context) {
		// Mock authenticated user in context
		user := &entity.AdminUser{
			ID:    1,
			Email: "test@example.com",
			Name:  "Test User",
		}
		c.Set("user", user)
		c.Next()
	}, handler.Me)

	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		expectedStatus int
		expectUser     bool
	}{
		{
			name: "Valid authenticated user",
			setupContext: func(c *gin.Context) {
				user := &entity.AdminUser{
					ID:    1,
					Email: "test@example.com",
					Name:  "Test User",
				}
				c.Set("user", user)
			},
			expectedStatus: http.StatusOK,
			expectUser:     true,
		},
		{
			name: "No user in context",
			setupContext: func(c *gin.Context) {
				// Don't set user in context
			},
			expectedStatus: http.StatusUnauthorized,
			expectUser:     false,
		},
		{
			name: "Invalid user type in context",
			setupContext: func(c *gin.Context) {
				c.Set("user", "invalid_user_type")
			},
			expectedStatus: http.StatusUnauthorized,
			expectUser:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create new router for each test
			router := gin.New()
			router.GET("/me", func(c *gin.Context) {
				tt.setupContext(c)
				c.Next()
			}, handler.Me)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/me", nil)
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response body
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if tt.expectUser {
				if _, exists := response["user"]; !exists {
					t.Error("Expected user in response")
				}
				
				// Verify user data
				if user, ok := response["user"].(map[string]interface{}); ok {
					if email, exists := user["email"]; !exists || email != "test@example.com" {
						t.Error("Expected correct user email in response")
					}
					if name, exists := user["name"]; !exists || name != "Test User" {
						t.Error("Expected correct user name in response")
					}
					// Password hash should not be included
					if _, exists := user["password_hash"]; exists {
						t.Error("Password hash should not be included in response")
					}
				} else {
					t.Error("User should be an object")
				}
			} else {
				if _, exists := response["error"]; !exists {
					t.Error("Expected error in response")
				}
			}
		})
	}
}

func TestHandler_LoginInvalidJSON(t *testing.T) {
	handler, router := setupTestHandler()
	router.POST("/login", handler.Login)

	// Test invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid JSON, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if _, exists := response["error"]; !exists {
		t.Error("Expected error in response for invalid JSON")
	}
}