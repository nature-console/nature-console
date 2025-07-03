package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/internal/utils"
	authUC "github.com/nature-console/backend/internal/usecase/auth"
	"github.com/nature-console/backend/test/mocks"
)

func TestSetupMiddlewares(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Create a new router
	router := gin.New()
	
	// Apply middlewares
	SetupMiddlewares(router)
	
	// Add a test endpoint
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	// Test that CORS middleware is applied
	t.Run("CORS middleware applied", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Check CORS headers are present
		allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
		if allowOrigin != "http://localhost:3000" {
			t.Errorf("Expected CORS header, got: %s", allowOrigin)
		}
	})

	// Test that logging middleware is applied
	t.Run("Logging middleware applied", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		// If logging middleware is working, the request should be processed normally
		// The actual logging output is tested in logging_test.go
	})

	// Test OPTIONS request (handled by CORS)
	t.Run("OPTIONS request handled", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "GET")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// OPTIONS should be handled by CORS middleware
		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status %d for OPTIONS, got %d", http.StatusNoContent, w.Code)
		}
	})
}

func TestNewAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Create auth use case
	mockRepo := mocks.NewMockAdminUserRepository()
	authUseCase := authUC.NewUseCase(mockRepo)
	
	// Create auth middleware using the helper function
	middleware := NewAuthMiddleware(authUseCase)
	
	if middleware == nil {
		t.Fatal("NewAuthMiddleware should return a non-nil middleware function")
	}

	// Test the middleware function
	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "protected"})
	})

	t.Run("Missing token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Invalid token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: "invalid.token",
		})
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestMiddleware_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Create router with all middlewares
	router := gin.New()
	SetupMiddlewares(router)
	
	// Add auth middleware
	mockRepo := mocks.NewMockAdminUserRepository()
	authUseCase := authUC.NewUseCase(mockRepo)
	authMiddleware := NewAuthMiddleware(authUseCase)
	
	// Protected route
	protectedGroup := router.Group("/api")
	protectedGroup.Use(authMiddleware)
	protectedGroup.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "protected"})
	})
	
	// Public route
	router.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "public"})
	})

	t.Run("Public route accessible", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/public", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		// CORS headers should be present
		allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
		if allowOrigin != "http://localhost:3000" {
			t.Error("CORS headers should be applied to public routes")
		}
	})

	t.Run("Protected route requires auth", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/protected", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Protected route with valid token", func(t *testing.T) {
		// Create test user
		user := &entity.AdminUser{
			ID:    1,
			Email: "test@example.com",
			Name:  "Test User",
		}
		mockRepo.Create(nil, user)
		
		// Generate token
		token, err := utils.GenerateToken(user.ID, user.Email)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		req := httptest.NewRequest("GET", "/api/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: token,
		})
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestMiddleware_ErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Create router with middlewares
	router := gin.New()
	
	// Add recovery middleware first
	router.Use(gin.Recovery())
	SetupMiddlewares(router)
	
	// Add a route that causes an error
	router.GET("/error", func(c *gin.Context) {
		panic("test panic")
	})

	t.Run("Middleware handles errors gracefully", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/error", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)

		// Should recover from panic and return 500
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d after panic recovery, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestMiddleware_MethodNotAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Create router with middlewares
	router := gin.New()
	SetupMiddlewares(router)
	
	// Only allow GET
	router.GET("/only-get", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	t.Run("Method not allowed", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/only-get", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 405 Method Not Allowed
		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}

		// CORS headers should still be applied
		allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
		if allowOrigin != "http://localhost:3000" {
			t.Error("CORS headers should be applied even for method not allowed")
		}
	})
}

func TestMiddleware_ChainOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	var executionOrder []string
	
	// Create router
	router := gin.New()
	
	// Add custom middleware to track execution order
	router.Use(func(c *gin.Context) {
		executionOrder = append(executionOrder, "custom1")
		c.Next()
	})
	
	// Apply standard middlewares
	SetupMiddlewares(router)
	
	// Add another custom middleware
	router.Use(func(c *gin.Context) {
		executionOrder = append(executionOrder, "custom2")
		c.Next()
	})
	
	router.GET("/order-test", func(c *gin.Context) {
		executionOrder = append(executionOrder, "handler")
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	req := httptest.NewRequest("GET", "/order-test", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check execution order
	expectedOrder := []string{"custom1", "custom2", "handler"}
	if len(executionOrder) != len(expectedOrder) {
		t.Errorf("Expected %d middleware executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if i >= len(executionOrder) || executionOrder[i] != expected {
			t.Errorf("Expected execution order[%d] to be %s, got %v", i, expected, executionOrder)
		}
	}
}