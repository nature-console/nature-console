package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoggingMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr) // Restore original output
	
	// Setup router with logging middleware
	router := gin.New()
	router.Use(LoggingMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"message": "created"})
	})

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		checkLogFor    []string
	}{
		{
			name:           "GET request",
			method:         "GET",
			path:           "/test",
			expectedStatus: http.StatusOK,
			checkLogFor:    []string{"GET", "/test", "200"},
		},
		{
			name:           "POST request",
			method:         "POST", 
			path:           "/test",
			expectedStatus: http.StatusCreated,
			checkLogFor:    []string{"POST", "/test", "201"},
		},
		{
			name:           "Not found request",
			method:         "GET",
			path:           "/notfound",
			expectedStatus: http.StatusNotFound,
			checkLogFor:    []string{"GET", "/notfound", "404"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear log buffer
			logBuffer.Reset()
			
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check log output
			logOutput := logBuffer.String()
			for _, checkStr := range tt.checkLogFor {
				if !strings.Contains(logOutput, checkStr) {
					t.Errorf("Expected log to contain '%s', got: %s", checkStr, logOutput)
				}
			}
		})
	}
}

func TestLoggingMiddleware_LogFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)
	
	// Setup router
	router := gin.New()
	router.Use(LoggingMiddleware())
	router.GET("/format-test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "test"})
	})

	req := httptest.NewRequest("GET", "/format-test", nil)
	req.Header.Set("User-Agent", "test-agent")
	req.RemoteAddr = "127.0.0.1:12345"
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	logOutput := logBuffer.String()
	
	// Check that log contains expected components
	expectedComponents := []string{
		"GET",           // HTTP method
		"/format-test",  // Path
		"200",           // Status code
		"127.0.0.1",     // Client IP (part of RemoteAddr)
	}

	for _, component := range expectedComponents {
		if !strings.Contains(logOutput, component) {
			t.Errorf("Expected log to contain '%s', got: %s", component, logOutput)
		}
	}
}

func TestLoggingMiddleware_DifferentStatusCodes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)
	
	// Setup router with different endpoints returning different status codes
	router := gin.New()
	router.Use(LoggingMiddleware())
	
	router.GET("/ok", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/created", func(c *gin.Context) {
		c.Status(http.StatusCreated)
	})
	router.GET("/badrequest", func(c *gin.Context) {
		c.Status(http.StatusBadRequest)
	})
	router.GET("/internal-error", func(c *gin.Context) {
		c.Status(http.StatusInternalServerError)
	})

	testCases := []struct {
		path           string
		expectedStatus string
	}{
		{"/ok", "200"},
		{"/created", "201"},
		{"/badrequest", "400"},
		{"/internal-error", "500"},
	}

	for _, tc := range testCases {
		t.Run("Status_"+tc.expectedStatus, func(t *testing.T) {
			logBuffer.Reset()
			
			req := httptest.NewRequest("GET", tc.path, nil)
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)

			logOutput := logBuffer.String()
			if !strings.Contains(logOutput, tc.expectedStatus) {
				t.Errorf("Expected log to contain status '%s', got: %s", tc.expectedStatus, logOutput)
			}
		})
	}
}

func TestLoggingMiddleware_ClientIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)
	
	// Setup router
	router := gin.New()
	router.Use(LoggingMiddleware())
	router.GET("/ip-test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	tests := []struct {
		name       string
		remoteAddr string
		headers    map[string]string
		expectIP   string
	}{
		{
			name:       "Direct connection",
			remoteAddr: "192.168.1.100:12345",
			expectIP:   "192.168.1.100",
		},
		{
			name:       "Localhost connection",
			remoteAddr: "127.0.0.1:54321",
			expectIP:   "127.0.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logBuffer.Reset()
			
			req := httptest.NewRequest("GET", "/ip-test", nil)
			req.RemoteAddr = tt.remoteAddr
			
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			logOutput := logBuffer.String()
			if !strings.Contains(logOutput, tt.expectIP) {
				t.Errorf("Expected log to contain IP '%s', got: %s", tt.expectIP, logOutput)
			}
		})
	}
}

func TestLoggingMiddleware_Latency(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)
	
	// Setup router
	router := gin.New()
	router.Use(LoggingMiddleware())
	router.GET("/latency-test", func(c *gin.Context) {
		// Simulate some processing time
		// Note: we don't actually sleep in tests to keep them fast
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/latency-test", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)

	logOutput := logBuffer.String()
	
	// The log should contain some latency information
	// Since we're not sleeping, latency should be very small (microseconds or nanoseconds)
	// Just check that the log was generated (latency will be included in the format)
	if logOutput == "" {
		t.Error("Expected log output, got empty string")
	}
	
	// Check that standard components are present
	expectedComponents := []string{"GET", "/latency-test", "200"}
	for _, component := range expectedComponents {
		if !strings.Contains(logOutput, component) {
			t.Errorf("Expected log to contain '%s', got: %s", component, logOutput)
		}
	}
}

func TestLoggingMiddleware_MiddlewareOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)
	
	// Setup router with multiple middlewares
	router := gin.New()
	
	// Add a custom middleware before logging
	router.Use(func(c *gin.Context) {
		c.Header("X-Custom", "test")
		c.Next()
	})
	
	router.Use(LoggingMiddleware())
	
	// Add a custom middleware after logging
	router.Use(func(c *gin.Context) {
		c.Header("X-After-Log", "test")
		c.Next()
	})
	
	router.GET("/order-test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/order-test", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)

	// Check that request was processed successfully
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check that headers from both middlewares are present
	if w.Header().Get("X-Custom") != "test" {
		t.Error("Expected X-Custom header to be set")
	}
	if w.Header().Get("X-After-Log") != "test" {
		t.Error("Expected X-After-Log header to be set")
	}

	// Check that logging occurred
	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, "GET") || !strings.Contains(logOutput, "/order-test") {
		t.Errorf("Expected log output, got: %s", logOutput)
	}
}