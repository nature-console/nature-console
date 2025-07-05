package repository

import (
	"context"
	"testing"
	"github.com/nature-console/backend/internal/domain/entity"
)

type MockArticleRepository struct {
	articles []*entity.Article
}

func (m *MockArticleRepository) Create(ctx context.Context, article *entity.Article) error {
	article.ID = uint(len(m.articles) + 1)
	m.articles = append(m.articles, article)
	return nil
}

func (m *MockArticleRepository) GetByID(ctx context.Context, id uint) (*entity.Article, error) {
	for _, article := range m.articles {
		if article.ID == id {
			return article, nil
		}
	}
	return nil, nil
}

func (m *MockArticleRepository) GetAll(ctx context.Context) ([]*entity.Article, error) {
	return m.articles, nil
}

func (m *MockArticleRepository) GetPublished(ctx context.Context) ([]*entity.Article, error) {
	var published []*entity.Article
	for _, article := range m.articles {
		if article.Published {
			published = append(published, article)
		}
	}
	return published, nil
}

func (m *MockArticleRepository) Update(ctx context.Context, article *entity.Article) error {
	for i, a := range m.articles {
		if a.ID == article.ID {
			m.articles[i] = article
			return nil
		}
	}
	return nil
}

func (m *MockArticleRepository) Delete(ctx context.Context, id uint) error {
	for i, article := range m.articles {
		if article.ID == id {
			m.articles = append(m.articles[:i], m.articles[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MockArticleRepository) GetByAuthor(ctx context.Context, author string) ([]*entity.Article, error) {
	var authorArticles []*entity.Article
	for _, article := range m.articles {
		if article.Author == author {
			authorArticles = append(authorArticles, article)
		}
	}
	return authorArticles, nil
}

func TestMockArticleRepository(t *testing.T) {
	repo := &MockArticleRepository{
		articles: []*entity.Article{},
	}
	
	// Test Create
	article := &entity.Article{
		Title:     "Test Article",
		Content:   "Test content",
		Author:    "Test Author",
		Published: true,
	}
	
	err := repo.Create(context.Background(), article)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if article.ID != 1 {
		t.Errorf("Expected ID 1, got %d", article.ID)
	}
	
	// Test GetByID
	found, err := repo.GetByID(context.Background(), 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if found.Title != "Test Article" {
		t.Errorf("Expected title %s, got %s", "Test Article", found.Title)
	}
}