package article

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/internal/usecase/article"
)

type mockUseCase struct {
	articles []*entity.Article
	nextID   uint
}

func (m *mockUseCase) CreateArticle(ctx context.Context, title, content, author string) (*entity.Article, error) {
	if title == "" || content == "" || author == "" {
		return nil, errors.New("validation error")
	}
	
	m.nextID++
	article := &entity.Article{
		ID:      m.nextID,
		Title:   title,
		Content: content,
		Author:  author,
	}
	m.articles = append(m.articles, article)
	return article, nil
}

func (m *mockUseCase) GetArticle(ctx context.Context, id uint) (*entity.Article, error) {
	for _, article := range m.articles {
		if article.ID == id {
			return article, nil
		}
	}
	return nil, errors.New("article not found")
}

func (m *mockUseCase) GetAllArticles(ctx context.Context) ([]*entity.Article, error) {
	return m.articles, nil
}

func (m *mockUseCase) GetPublishedArticles(ctx context.Context) ([]*entity.Article, error) {
	var published []*entity.Article
	for _, article := range m.articles {
		if article.Published {
			published = append(published, article)
		}
	}
	return published, nil
}

func (m *mockUseCase) UpdateArticle(ctx context.Context, id uint, title, content, author string, published bool) (*entity.Article, error) {
	for _, article := range m.articles {
		if article.ID == id {
			if title != "" {
				article.Title = title
			}
			if content != "" {
				article.Content = content
			}
			if author != "" {
				article.Author = author
			}
			article.Published = published
			return article, nil
		}
	}
	return nil, errors.New("article not found")
}

func (m *mockUseCase) DeleteArticle(ctx context.Context, id uint) error {
	for i, article := range m.articles {
		if article.ID == id {
			m.articles = append(m.articles[:i], m.articles[i+1:]...)
			return nil
		}
	}
	return errors.New("article not found")
}

func (m *mockUseCase) GetArticlesByAuthor(ctx context.Context, author string) ([]*entity.Article, error) {
	var authorArticles []*entity.Article
	for _, article := range m.articles {
		if article.Author == author {
			authorArticles = append(authorArticles, article)
		}
	}
	return authorArticles, nil
}

func (m *mockUseCase) PublishArticle(ctx context.Context, id uint) (*entity.Article, error) {
	for _, article := range m.articles {
		if article.ID == id {
			article.Published = true
			return article, nil
		}
	}
	return nil, errors.New("article not found")
}

func (m *mockUseCase) UnpublishArticle(ctx context.Context, id uint) (*entity.Article, error) {
	for _, article := range m.articles {
		if article.ID == id {
			article.Published = false
			return article, nil
		}
	}
	return nil, errors.New("article not found")
}

func TestHandler_CreateArticle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockUC := &mockUseCase{
		articles: []*entity.Article{},
		nextID:   0,
	}
	
	// Create a fake UseCase that matches the interface
	realUC := article.NewUseCase(nil) // We'll use mock instead
	handler := NewHandler(realUC)
	
	// For testing, we'll replace the useCase with our mock
	handler.useCase = (*article.UseCase)(mockUC)
	
	router := gin.New()
	router.POST("/articles", handler.CreateArticle)
	
	// Test successful creation
	reqBody := CreateArticleRequest{
		Title:   "Test Article",
		Content: "Test Content",
		Author:  "Test Author",
	}
	
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
	
	// Test validation error
	reqBody = CreateArticleRequest{
		Title:   "",
		Content: "Test Content",
		Author:  "Test Author",
	}
	
	jsonBody, _ = json.Marshal(reqBody)
	req = httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_GetArticle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockUC := &mockUseCase{
		articles: []*entity.Article{
			{ID: 1, Title: "Test Article", Content: "Test Content", Author: "Test Author"},
		},
		nextID: 1,
	}
	
	realUC := article.NewUseCase(nil)
	handler := NewHandler(realUC)
	handler.useCase = (*article.UseCase)(mockUC)
	
	router := gin.New()
	router.GET("/articles/:id", handler.GetArticle)
	
	// Test successful retrieval
	req := httptest.NewRequest(http.MethodGet, "/articles/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	// Test invalid ID
	req = httptest.NewRequest(http.MethodGet, "/articles/invalid", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}