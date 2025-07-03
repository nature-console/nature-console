package article

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/internal/usecase/article"
	"github.com/nature-console/backend/test/mocks"
)

func TestHandler_CreateArticle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Create a mock repository and use it with the real usecase
	mockRepo := mocks.NewMockArticleRepository()
	useCase := article.NewUseCase(mockRepo)
	handler := NewHandler(useCase)
	
	router := gin.New()
	router.POST("/articles", handler.CreateArticle)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Valid article creation",
			requestBody: map[string]interface{}{
				"title":   "Test Article",
				"content": "Test Content",
				"author":  "Test Author",
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "Missing title",
			requestBody: map[string]interface{}{
				"content": "Test Content",
				"author":  "Test Author",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Missing content",
			requestBody: map[string]interface{}{
				"title":  "Test Article",
				"author": "Test Author",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Missing author",
			requestBody: map[string]interface{}{
				"title":   "Test Article",
				"content": "Test Content",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Errorf("Failed to unmarshal response: %v", err)
			}

			if tt.expectError {
				if _, exists := response["error"]; !exists {
					t.Error("Expected error in response")
				}
			} else {
				if _, exists := response["error"]; exists {
					t.Errorf("Unexpected error in response: %v", response["error"])
				}
				if _, exists := response["id"]; !exists {
					t.Error("Expected article ID in response")
				}
			}
		})
	}
}

func TestHandler_CreateArticleInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockRepo := mocks.NewMockArticleRepository()
	useCase := article.NewUseCase(mockRepo)
	handler := NewHandler(useCase)

	router := gin.New()
	router.POST("/articles", handler.CreateArticle)

	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_GetArticle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Create a mock repository and use it with the real usecase
	mockRepo := mocks.NewMockArticleRepository()
	useCase := article.NewUseCase(mockRepo)
	handler := NewHandler(useCase)
	
	// Add test article to mock repository
	testArticle := &entity.Article{
		Title:     "Test Article",
		Content:   "Test Content",
		Author:    "Test Author",
		Published: true,
	}
	mockRepo.Create(context.Background(), testArticle)
	
	router := gin.New()
	router.GET("/articles/:id", handler.GetArticle)

	tests := []struct {
		name           string
		articleID      string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid article ID",
			articleID:      "1",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Invalid article ID",
			articleID:      "invalid",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Non-existent article",
			articleID:      "999",
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/articles/"+tt.articleID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Errorf("Failed to unmarshal response: %v", err)
			}

			if tt.expectError {
				if _, exists := response["error"]; !exists {
					t.Error("Expected error in response")
				}
			} else {
				if _, exists := response["error"]; exists {
					t.Errorf("Unexpected error in response: %v", response["error"])
				}
				if _, exists := response["id"]; !exists {
					t.Error("Expected article ID in response")
				}
			}
		})
	}
}