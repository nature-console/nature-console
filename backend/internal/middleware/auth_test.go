package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/internal/utils"
	authUC "github.com/nature-console/backend/internal/usecase/auth"
	"github.com/nature-console/backend/test/mocks"
)

func setupTestAuthMiddleware() (*authUC.UseCase, gin.HandlerFunc) {
	gin.SetMode(gin.TestMode)
	mockRepo := mocks.NewMockAdminUserRepository()
	useCase := authUC.NewUseCase(mockRepo)
	middleware := AuthMiddleware(useCase)
	return useCase, middleware
}

func createTestUserForMiddleware(useCase *authUC.UseCase) (*entity.AdminUser, string) {
	// Get the mock repository from the use case (we'll need to create a user)
	mockRepo := mocks.NewMockAdminUserRepository()
	useCase = authUC.NewUseCase(mockRepo)

	user := &entity.AdminUser{
		Email:        "middleware@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Middleware Test User",
	}
	
	mockRepo.Create(nil, user)
	
	// Generate a token for this user
	token, _ := utils.GenerateToken(user.ID, user.Email)
	return user, token
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	useCase, middleware := setupTestAuthMiddleware()
	
	// Create test user and token
	user, token := createTestUserForMiddleware(useCase)
	
	// Create a fresh use case with the user
	mockRepo := mocks.NewMockAdminUserRepository()
	mockRepo.Create(nil, user)
	useCase = authUC.NewUseCase(mockRepo)
	middleware = AuthMiddleware(useCase)
	
	// Setup router
	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		userFromContext, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": userFromContext})
	})
	
	// Create request with valid token cookie
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if _, exists := response["user"]; !exists {
		t.Error("Expected user in response")
	}
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	_, middleware := setupTestAuthMiddleware()
	
	// Setup router
	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	// Create request without token cookie
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if error_msg, exists := response["error"]; !exists || error_msg != "Authentication required" {
		t.Error("Expected 'Authentication required' error message")
	}
}

func TestAuthMiddleware_EmptyToken(t *testing.T) {
	_, middleware := setupTestAuthMiddleware()
	
	// Setup router
	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	// Create request with empty token cookie
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "",
	})
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	_, middleware := setupTestAuthMiddleware()
	
	// Setup router
	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	// Create request with invalid token cookie
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "invalid.token.here",
	})
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if error_msg, exists := response["error"]; !exists || error_msg != "Invalid token" {
		t.Error("Expected 'Invalid token' error message")
	}
}

func TestAuthMiddleware_UserNotFound(t *testing.T) {
	_, middleware := setupTestAuthMiddleware()
	
	// Create token for non-existent user
	token, err := utils.GenerateToken(999, "nonexistent@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// Setup router
	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	// Create request with token for non-existent user
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestValidateTokenMiddleware_ValidToken(t *testing.T) {
	// Generate a valid token
	userID := uint(123)
	email := "test@example.com"
	token, err := utils.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ValidateTokenMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		userIDFromContext, userIDExists := c.Get("userID")
		emailFromContext, emailExists := c.Get("email")
		
		if !userIDExists || !emailExists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Claims not found in context"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"userID": userIDFromContext,
			"email":  emailFromContext,
		})
	})
	
	// Create request with valid token cookie
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if responseUserID := response["userID"].(float64); uint(responseUserID) != userID {
		t.Errorf("Expected userID %d, got %v", userID, responseUserID)
	}
	
	if responseEmail := response["email"].(string); responseEmail != email {
		t.Errorf("Expected email %s, got %s", email, responseEmail)
	}
}

func TestValidateTokenMiddleware_MissingToken(t *testing.T) {
	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ValidateTokenMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	// Create request without token cookie
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if error_msg, exists := response["error"]; !exists || error_msg != "Authentication required" {
		t.Error("Expected 'Authentication required' error message")
	}
}

func TestValidateTokenMiddleware_InvalidToken(t *testing.T) {
	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ValidateTokenMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	// Create request with invalid token cookie
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "invalid.token.here",
	})
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if error_msg, exists := response["error"]; !exists || error_msg != "Invalid token" {
		t.Error("Expected 'Invalid token' error message")
	}
}

func TestValidateTokenMiddleware_EmptyToken(t *testing.T) {
	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ValidateTokenMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	// Create request with empty token cookie
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "",
	})
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestMiddleware_AbortBehavior(t *testing.T) {
	// Test that middleware properly aborts the request chain
	_, middleware := setupTestAuthMiddleware()
	
	handlerCalled := false
	
	// Setup router
	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		handlerCalled = true
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	// Create request without token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	if handlerCalled {
		t.Error("Handler should not be called when middleware aborts")
	}
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestMiddleware_NextBehavior(t *testing.T) {
	// Test that middleware calls Next() for valid requests
	useCase, middleware := setupTestAuthMiddleware()
	
	// Create test user and token
	user, token := createTestUserForMiddleware(useCase)
	
	// Create a fresh use case with the user
	mockRepo := mocks.NewMockAdminUserRepository()
	mockRepo.Create(nil, user)
	useCase = authUC.NewUseCase(mockRepo)
	middleware = AuthMiddleware(useCase)
	
	handlerCalled := false
	
	// Setup router
	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		handlerCalled = true
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	// Create request with valid token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if !handlerCalled {
		t.Error("Handler should be called when middleware succeeds")
	}
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestMiddleware_ContextValues(t *testing.T) {
	// Test ValidateTokenMiddleware sets correct context values
	userID := uint(456)
	email := "context@example.com"
	token, err := utils.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	var contextUserID uint
	var contextEmail string
	var userIDExists, emailExists bool
	
	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ValidateTokenMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		userIDInterface, userIDExists := c.Get("userID")
		emailInterface, emailExists := c.Get("email")
		
		if userIDExists {
			contextUserID = userIDInterface.(uint)
		}
		if emailExists {
			contextEmail = emailInterface.(string)
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	// Create request with valid token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		return
	}
	
	if !userIDExists {
		t.Error("userID should be set in context")
	}
	if !emailExists {
		t.Error("email should be set in context")
	}
	
	if userIDExists && contextUserID != userID {
		t.Errorf("Expected userID %d in context, got %d", userID, contextUserID)
	}
	if emailExists && contextEmail != email {
		t.Errorf("Expected email %s in context, got %s", email, contextEmail)
	}
}