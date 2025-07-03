package mocks

import (
	"context"
	"errors"

	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/internal/domain/repository"
)

type MockArticleRepository struct {
	articles []*entity.Article
	nextID   uint
}

func NewMockArticleRepository() repository.ArticleRepository {
	return &MockArticleRepository{
		articles: make([]*entity.Article, 0),
		nextID:   1,
	}
}

func (m *MockArticleRepository) Create(ctx context.Context, article *entity.Article) error {
	if article == nil {
		return errors.New("article cannot be nil")
	}
	
	article.ID = m.nextID
	m.nextID++
	
	m.articles = append(m.articles, article)
	return nil
}

func (m *MockArticleRepository) GetByID(ctx context.Context, id uint) (*entity.Article, error) {
	for _, article := range m.articles {
		if article.ID == id {
			return article, nil
		}
	}
	return nil, errors.New("article not found")
}

func (m *MockArticleRepository) GetAll(ctx context.Context) ([]*entity.Article, error) {
	result := make([]*entity.Article, len(m.articles))
	copy(result, m.articles)
	return result, nil
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
	if article == nil {
		return errors.New("article cannot be nil")
	}
	
	for i, a := range m.articles {
		if a.ID == article.ID {
			m.articles[i] = article
			return nil
		}
	}
	return errors.New("article not found")
}

func (m *MockArticleRepository) Delete(ctx context.Context, id uint) error {
	for i, article := range m.articles {
		if article.ID == id {
			m.articles = append(m.articles[:i], m.articles[i+1:]...)
			return nil
		}
	}
	return errors.New("article not found")
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