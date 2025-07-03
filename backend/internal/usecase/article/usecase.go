package article

import (
	"context"
	"errors"
	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/internal/domain/repository"
)

type UseCase struct {
	articleRepo repository.ArticleRepository
}

func NewUseCase(articleRepo repository.ArticleRepository) *UseCase {
	return &UseCase{
		articleRepo: articleRepo,
	}
}

func (u *UseCase) CreateArticle(ctx context.Context, title, content, author string) (*entity.Article, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if content == "" {
		return nil, errors.New("content is required")
	}
	if author == "" {
		return nil, errors.New("author is required")
	}
	
	article := &entity.Article{
		Title:     title,
		Content:   content,
		Author:    author,
		Published: false,
	}
	
	err := u.articleRepo.Create(ctx, article)
	if err != nil {
		return nil, err
	}
	
	return article, nil
}

func (u *UseCase) GetArticle(ctx context.Context, id uint) (*entity.Article, error) {
	if id == 0 {
		return nil, errors.New("invalid article ID")
	}
	
	return u.articleRepo.GetByID(ctx, id)
}

func (u *UseCase) GetAllArticles(ctx context.Context) ([]*entity.Article, error) {
	return u.articleRepo.GetAll(ctx)
}

func (u *UseCase) GetPublishedArticles(ctx context.Context) ([]*entity.Article, error) {
	return u.articleRepo.GetPublished(ctx)
}

func (u *UseCase) UpdateArticle(ctx context.Context, id uint, title, content, author string, published bool) (*entity.Article, error) {
	if id == 0 {
		return nil, errors.New("invalid article ID")
	}
	
	article, err := u.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
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
	
	err = u.articleRepo.Update(ctx, article)
	if err != nil {
		return nil, err
	}
	
	return article, nil
}

func (u *UseCase) DeleteArticle(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid article ID")
	}
	
	return u.articleRepo.Delete(ctx, id)
}

func (u *UseCase) GetArticlesByAuthor(ctx context.Context, author string) ([]*entity.Article, error) {
	if author == "" {
		return nil, errors.New("author is required")
	}
	
	return u.articleRepo.GetByAuthor(ctx, author)
}

func (u *UseCase) PublishArticle(ctx context.Context, id uint) (*entity.Article, error) {
	if id == 0 {
		return nil, errors.New("invalid article ID")
	}
	
	article, err := u.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	article.Published = true
	
	err = u.articleRepo.Update(ctx, article)
	if err != nil {
		return nil, err
	}
	
	return article, nil
}

func (u *UseCase) UnpublishArticle(ctx context.Context, id uint) (*entity.Article, error) {
	if id == 0 {
		return nil, errors.New("invalid article ID")
	}
	
	article, err := u.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	article.Published = false
	
	err = u.articleRepo.Update(ctx, article)
	if err != nil {
		return nil, err
	}
	
	return article, nil
}