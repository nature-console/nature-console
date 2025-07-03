package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nature-console/backend/internal/domain/entity"
	articleUC "github.com/nature-console/backend/internal/usecase/article"
	"github.com/nature-console/backend/test/mocks"
)

func setupTestHandler() (*Handler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockRepo := mocks.NewMockArticleRepository()
	useCase := articleUC.NewUseCase(mockRepo)
	handler := NewHandler(useCase)
	
	router := gin.New()
	return handler, router
}

func createTestArticle(title, content, author string, published bool) *entity.Article {
	article := &entity.Article{
		Title:     title,
		Content:   content,
		Author:    author,
		Published: published,
	}
	return article
}

func TestHandler_GetDashboard(t *testing.T) {
	handler, router := setupTestHandler()
	router.GET("/dashboard", handler.GetDashboard)

	// Setup test data
	mockRepo := mocks.NewMockArticleRepository()
	useCase := articleUC.NewUseCase(mockRepo)
	handler = NewHandler(useCase)
	router.GET("/dashboard", handler.GetDashboard)

	// Create test articles
	article1 := createTestArticle("Published Article 1", "Content 1", "Author 1", true)
	article2 := createTestArticle("Published Article 2", "Content 2", "Author 2", true)
	article3 := createTestArticle("Draft Article 1", "Content 3", "Author 3", false)
	
	mockRepo.Create(nil, article1)
	mockRepo.Create(nil, article2)
	mockRepo.Create(nil, article3)

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
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

	// Check stats
	stats, exists := response["stats"].(map[string]interface{})
	if !exists {
		t.Fatal("Expected stats in response")
	}

	if totalArticles := stats["total_articles"].(float64); totalArticles != 3 {
		t.Errorf("Expected 3 total articles, got %v", totalArticles)
	}
	if publishedArticles := stats["published_articles"].(float64); publishedArticles != 2 {
		t.Errorf("Expected 2 published articles, got %v", publishedArticles)
	}
	if draftArticles := stats["draft_articles"].(float64); draftArticles != 1 {
		t.Errorf("Expected 1 draft article, got %v", draftArticles)
	}

	// Check recent articles
	recentArticles, exists := response["recent_articles"].([]interface{})
	if !exists {
		t.Fatal("Expected recent_articles in response")
	}
	if len(recentArticles) != 3 {
		t.Errorf("Expected 3 recent articles, got %d", len(recentArticles))
	}
}

func TestHandler_GetAllArticles(t *testing.T) {
	handler, router := setupTestHandler()
	router.GET("/articles", handler.GetAllArticles)

	// Setup test data
	mockRepo := mocks.NewMockArticleRepository()
	useCase := articleUC.NewUseCase(mockRepo)
	handler = NewHandler(useCase)
	router.GET("/articles", handler.GetAllArticles)

	article := createTestArticle( "Article 1", "Content 1", "Author 1", true)
	article := createTestArticle( "Article 2", "Content 2", "Author 2", false)

	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var articles []entity.Article
	err := json.Unmarshal(w.Body.Bytes(), &articles)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(articles) != 2 {
		t.Errorf("Expected 2 articles, got %d", len(articles))
	}
}

func TestHandler_CreateArticle(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateArticleRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Valid draft article",
			requestBody: CreateArticleRequest{
				Title:     "Test Article",
				Content:   "Test Content",
				Author:    "Test Author",
				Published: false,
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "Valid published article",
			requestBody: CreateArticleRequest{
				Title:     "Published Article",
				Content:   "Published Content",
				Author:    "Published Author",
				Published: true,
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "Missing title",
			requestBody: CreateArticleRequest{
				Content:   "Test Content",
				Author:    "Test Author",
				Published: false,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Missing content",
			requestBody: CreateArticleRequest{
				Title:     "Test Article",
				Author:    "Test Author",
				Published: false,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Missing author",
			requestBody: CreateArticleRequest{
				Title:     "Test Article",
				Content:   "Test Content",
				Published: false,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, router := setupTestHandler()
			router.POST("/articles", handler.CreateArticle)

			requestBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
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

func TestHandler_UpdateArticle(t *testing.T) {
	handler, router := setupTestHandler()
	
	// Setup test data
	mockRepo := mocks.NewMockArticleRepository()
	useCase := articleUC.NewUseCase(mockRepo)
	handler = NewHandler(useCase)

	article := article := createTestArticle( "Original Title", "Original Content", "Original Author", false)
	
	router.PUT("/articles/:id", handler.UpdateArticle)

	tests := []struct {
		name           string
		articleID      string
		requestBody    UpdateArticleRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name:      "Valid update",
			articleID: strconv.Itoa(int(article.ID)),
			requestBody: UpdateArticleRequest{
				Title:     "Updated Title",
				Content:   "Updated Content",
				Author:    "Updated Author",
				Published: true,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:      "Invalid article ID",
			articleID: "invalid",
			requestBody: UpdateArticleRequest{
				Title:   "Updated Title",
				Content: "Updated Content",
				Author:  "Updated Author",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "Non-existent article",
			articleID: "999",
			requestBody: UpdateArticleRequest{
				Title:   "Updated Title",
				Content: "Updated Content",
				Author:  "Updated Author",
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPut, "/articles/"+tt.articleID, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if tt.expectError {
				if _, exists := response["error"]; !exists {
					t.Error("Expected error in response")
				}
			} else {
				if _, exists := response["error"]; exists {
					t.Errorf("Unexpected error in response: %v", response["error"])
				}
			}
		})
	}
}

func TestHandler_DeleteArticle(t *testing.T) {
	handler, router := setupTestHandler()
	
	// Setup test data
	mockRepo := mocks.NewMockArticleRepository()
	useCase := articleUC.NewUseCase(mockRepo)
	handler = NewHandler(useCase)

	article := article := createTestArticle( "To Delete", "Content", "Author", false)
	
	router.DELETE("/articles/:id", handler.DeleteArticle)

	tests := []struct {
		name           string
		articleID      string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid deletion",
			articleID:      strconv.Itoa(int(article.ID)),
			expectedStatus: http.StatusNoContent,
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
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/articles/"+tt.articleID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectError && w.Code != http.StatusNoContent {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err == nil {
					if _, exists := response["error"]; !exists {
						t.Error("Expected error in response")
					}
				}
			}
		})
	}
}

func TestHandler_GetArticle(t *testing.T) {
	handler, router := setupTestHandler()
	
	// Setup test data
	mockRepo := mocks.NewMockArticleRepository()
	useCase := articleUC.NewUseCase(mockRepo)
	handler = NewHandler(useCase)

	article := article := createTestArticle( "Test Article", "Test Content", "Test Author", true)
	
	router.GET("/articles/:id", handler.GetArticle)

	tests := []struct {
		name           string
		articleID      string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid article ID",
			articleID:      strconv.Itoa(int(article.ID)),
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
				t.Fatalf("Failed to unmarshal response: %v", err)
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
	handler, router := setupTestHandler()
	router.POST("/articles", handler.CreateArticle)

	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBufferString("{invalid json"))
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

func TestDashboardStats_Structure(t *testing.T) {
	stats := DashboardStats{
		TotalArticles:     10,
		PublishedArticles: 7,
		DraftArticles:     3,
	}

	if stats.TotalArticles != 10 {
		t.Errorf("Expected TotalArticles 10, got %d", stats.TotalArticles)
	}
	if stats.PublishedArticles != 7 {
		t.Errorf("Expected PublishedArticles 7, got %d", stats.PublishedArticles)
	}
	if stats.DraftArticles != 3 {
		t.Errorf("Expected DraftArticles 3, got %d", stats.DraftArticles)
	}

	// Test JSON serialization
	jsonBytes, err := json.Marshal(stats)
	if err != nil {
		t.Fatalf("Failed to marshal DashboardStats to JSON: %v", err)
	}

	var unmarshaled DashboardStats
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal DashboardStats from JSON: %v", err)
	}

	if unmarshaled != stats {
		t.Error("DashboardStats did not survive JSON round trip")
	}
}

func TestHandler_GetDashboardWithManyArticles(t *testing.T) {
	handler, router := setupTestHandler()
	
	// Setup test data with more than 5 articles
	mockRepo := mocks.NewMockArticleRepository()
	useCase := articleUC.NewUseCase(mockRepo)
	handler = NewHandler(useCase)
	router.GET("/dashboard", handler.GetDashboard)

	// Create 7 articles
	for i := 0; i < 7; i++ {
		published := i%2 == 0 // Alternate between published and draft
		article := createTestArticle( "Article "+strconv.Itoa(i), "Content "+strconv.Itoa(i), "Author "+strconv.Itoa(i), published)
	}

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
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

	// Check that recent articles are limited to 5
	recentArticles, exists := response["recent_articles"].([]interface{})
	if !exists {
		t.Fatal("Expected recent_articles in response")
	}
	if len(recentArticles) != 5 {
		t.Errorf("Expected 5 recent articles, got %d", len(recentArticles))
	}

	// Check stats
	stats, exists := response["stats"].(map[string]interface{})
	if !exists {
		t.Fatal("Expected stats in response")
	}
	if totalArticles := stats["total_articles"].(float64); totalArticles != 7 {
		t.Errorf("Expected 7 total articles, got %v", totalArticles)
	}
}