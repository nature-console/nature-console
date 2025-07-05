package article

import (
	"context"
	"errors"
	"testing"
	"github.com/nature-console/backend/internal/domain/entity"
)

type mockArticleRepository struct {
	articles []*entity.Article
	nextID   uint
}

func (m *mockArticleRepository) Create(ctx context.Context, article *entity.Article) error {
	m.nextID++
	article.ID = m.nextID
	m.articles = append(m.articles, article)
	return nil
}

func (m *mockArticleRepository) GetByID(ctx context.Context, id uint) (*entity.Article, error) {
	for _, article := range m.articles {
		if article.ID == id {
			return article, nil
		}
	}
	return nil, errors.New("article not found")
}

func (m *mockArticleRepository) GetAll(ctx context.Context) ([]*entity.Article, error) {
	return m.articles, nil
}

func (m *mockArticleRepository) GetPublished(ctx context.Context) ([]*entity.Article, error) {
	var published []*entity.Article
	for _, article := range m.articles {
		if article.Published {
			published = append(published, article)
		}
	}
	return published, nil
}

func (m *mockArticleRepository) Update(ctx context.Context, article *entity.Article) error {
	for i, a := range m.articles {
		if a.ID == article.ID {
			m.articles[i] = article
			return nil
		}
	}
	return errors.New("article not found")
}

func (m *mockArticleRepository) Delete(ctx context.Context, id uint) error {
	for i, article := range m.articles {
		if article.ID == id {
			m.articles = append(m.articles[:i], m.articles[i+1:]...)
			return nil
		}
	}
	return errors.New("article not found")
}

func (m *mockArticleRepository) GetByAuthor(ctx context.Context, author string) ([]*entity.Article, error) {
	var authorArticles []*entity.Article
	for _, article := range m.articles {
		if article.Author == author {
			authorArticles = append(authorArticles, article)
		}
	}
	return authorArticles, nil
}

func TestUseCase_CreateArticle(t *testing.T) {
	repo := &mockArticleRepository{
		articles: []*entity.Article{},
		nextID:   0,
	}
	
	useCase := NewUseCase(repo)
	
	// Test successful creation
	article, err := useCase.CreateArticle(context.Background(), "Test Title", "Test Content", "Test Author")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if article.ID != 1 {
		t.Errorf("Expected ID 1, got %d", article.ID)
	}
	
	if article.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", article.Title)
	}
	
	if article.Published {
		t.Error("Expected published to be false for new article")
	}
	
	// Test validation errors
	_, err = useCase.CreateArticle(context.Background(), "", "Test Content", "Test Author")
	if err == nil {
		t.Error("Expected error for empty title")
	}
	
	_, err = useCase.CreateArticle(context.Background(), "Test Title", "", "Test Author")
	if err == nil {
		t.Error("Expected error for empty content")
	}
	
	_, err = useCase.CreateArticle(context.Background(), "Test Title", "Test Content", "")
	if err == nil {
		t.Error("Expected error for empty author")
	}
}

func TestUseCase_GetArticle(t *testing.T) {
	repo := &mockArticleRepository{
		articles: []*entity.Article{
			{ID: 1, Title: "Test Article", Content: "Test Content", Author: "Test Author"},
		},
		nextID: 1,
	}
	
	useCase := NewUseCase(repo)
	
	// Test successful retrieval
	article, err := useCase.GetArticle(context.Background(), 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if article.Title != "Test Article" {
		t.Errorf("Expected title 'Test Article', got '%s'", article.Title)
	}
	
	// Test validation error
	_, err = useCase.GetArticle(context.Background(), 0)
	if err == nil {
		t.Error("Expected error for invalid ID")
	}
}

func TestUseCase_PublishArticle(t *testing.T) {
	repo := &mockArticleRepository{
		articles: []*entity.Article{
			{ID: 1, Title: "Test Article", Content: "Test Content", Author: "Test Author", Published: false},
		},
		nextID: 1,
	}
	
	useCase := NewUseCase(repo)
	
	// Test successful publish
	article, err := useCase.PublishArticle(context.Background(), 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if !article.Published {
		t.Error("Expected article to be published")
	}
	
	// Test validation error
	_, err = useCase.PublishArticle(context.Background(), 0)
	if err == nil {
		t.Error("Expected error for invalid ID")
	}
}