package repository

import (
	"context"
	"github.com/nature-console/backend/internal/domain/entity"
)

type ArticleRepository interface {
	Create(ctx context.Context, article *entity.Article) error
	GetByID(ctx context.Context, id uint) (*entity.Article, error)
	GetAll(ctx context.Context) ([]*entity.Article, error)
	GetPublished(ctx context.Context) ([]*entity.Article, error)
	Update(ctx context.Context, article *entity.Article) error
	Delete(ctx context.Context, id uint) error
	GetByAuthor(ctx context.Context, author string) ([]*entity.Article, error)
}