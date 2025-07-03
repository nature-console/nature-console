package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Setup router with CORS middleware
	router := gin.New()
	router.Use(CORSMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	tests := []struct {
		name           string
		method         string
		origin         string
		expectedStatus int
		checkHeaders   map[string]string
	}{
		{
			name:           "GET request with allowed origin",
			method:         "GET",
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusOK,
			checkHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "http://localhost:3000",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			name:           "POST request with allowed origin",
			method:         "POST",
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusOK,
			checkHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "http://localhost:3000",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			name:           "OPTIONS preflight request",
			method:         "OPTIONS",
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusNoContent,
			checkHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "http://localhost:3000",
				"Access-Control-Allow-Methods":     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
				"Access-Control-Allow-Headers":     "Origin,Content-Type,Authorization",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			name:           "Request without origin",
			method:         "GET",
			origin:         "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check CORS headers
			for header, expectedValue := range tt.checkHeaders {
				actualValue := w.Header().Get(header)
				if actualValue != expectedValue {
					t.Errorf("Expected header %s: %s, got: %s", header, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestCORSMiddleware_AllowedMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Setup router
	router := gin.New()
	router.Use(CORSMiddleware())
	
	// Add handlers for all allowed methods
	router.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.POST("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.PUT("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.DELETE("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.PATCH("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	allowedMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}

	for _, method := range allowedMethods {
		t.Run("Method_"+method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/test", nil)
			req.Header.Set("Origin", "http://localhost:3000")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// All methods should be handled (either 200 OK or 204 No Content for OPTIONS)
			if w.Code != http.StatusOK && w.Code != http.StatusNoContent {
				t.Errorf("Method %s should be allowed, got status %d", method, w.Code)
			}
		})
	}
}

func TestCORSMiddleware_DisallowedOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Setup router
	router := gin.New()
	router.Use(CORSMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	// Test with disallowed origin
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://malicious-site.com")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Request might be blocked by CORS middleware
	// The exact behavior depends on the CORS configuration
	if w.Code != http.StatusOK && w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d or %d, got %d", http.StatusOK, http.StatusForbidden, w.Code)
	}

	// Check that disallowed origin is not in Access-Control-Allow-Origin
	allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin == "http://malicious-site.com" {
		t.Error("Disallowed origin should not be in Access-Control-Allow-Origin header")
	}
}

func TestCORSMiddleware_Headers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Setup router
	router := gin.New()
	router.Use(CORSMiddleware())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	// Test with various headers
	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check that Content-Length is exposed
	exposeHeaders := w.Header().Get("Access-Control-Expose-Headers")
	if exposeHeaders != "Content-Length" {
		t.Errorf("Expected Access-Control-Expose-Headers: Content-Length, got: %s", exposeHeaders)
	}
}

func TestCORSMiddleware_Credentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Setup router
	router := gin.New()
	router.Use(CORSMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check that credentials are allowed
	allowCredentials := w.Header().Get("Access-Control-Allow-Credentials")
	if allowCredentials != "true" {
		t.Errorf("Expected Access-Control-Allow-Credentials: true, got: %s", allowCredentials)
	}
}