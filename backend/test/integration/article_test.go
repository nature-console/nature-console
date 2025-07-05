package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/internal/handler/article"
	articleRepo "github.com/nature-console/backend/internal/repository/article"
	"github.com/nature-console/backend/internal/routes"
	articleUC "github.com/nature-console/backend/internal/usecase/article"
)

func setupTestDB() *gorm.DB {
	// Use a test database URL from environment or a default test database
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		// Skip integration tests if no test database is configured
		return nil
	}
	
	db, err := gorm.Open(postgres.Open(testDBURL), &gorm.Config{})
	if err != nil {
		return nil
	}
	
	// Auto migrate for testing
	db.AutoMigrate(&entity.Article{})
	
	return db
}

func cleanupTestDB(db *gorm.DB) {
	// Clean up test data
	db.Exec("TRUNCATE TABLE articles")
}

func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	
	// Initialize repositories
	articleRepository := articleRepo.NewArticleRepository(db)
	
	// Initialize use cases
	articleUseCase := articleUC.NewUseCase(articleRepository)
	
	// Initialize handlers
	articleHandler := article.NewHandler(articleUseCase)
	
	// Setup router
	r := gin.New()
	api := r.Group("/api/v1")
	routes.SetupArticleRoutes(api, articleHandler)
	
	return r
}

func TestArticleIntegration(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("Skipping integration test - no test database configured")
	}
	defer cleanupTestDB(db)
	
	router := setupRouter(db)
	
	t.Run("CreateArticle", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title":   "Integration Test Article",
			"content": "This is a test article content",
			"author":  "Test Author",
		}
		
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/articles", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
		}
		
		var response entity.Article
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}
		
		if response.Title != "Integration Test Article" {
			t.Errorf("Expected title 'Integration Test Article', got '%s'", response.Title)
		}
	})
	
	t.Run("GetAllArticles", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/articles", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
		
		var response []*entity.Article
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}
		
		if len(response) == 0 {
			t.Error("Expected at least one article")
		}
	})
	
	t.Run("GetArticle", func(t *testing.T) {
		// First create an article
		reqBody := map[string]interface{}{
			"title":   "Get Test Article",
			"content": "This is a get test article content",
			"author":  "Get Test Author",
		}
		
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/articles", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		var createdArticle entity.Article
		json.Unmarshal(w.Body.Bytes(), &createdArticle)
		
		// Now get the article
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/articles/%d", createdArticle.ID), nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
		
		var response entity.Article
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}
		
		if response.Title != "Get Test Article" {
			t.Errorf("Expected title 'Get Test Article', got '%s'", response.Title)
		}
	})
}