package article

import (
	"context"
	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) repository.ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(ctx context.Context, article *entity.Article) error {
	return r.db.WithContext(ctx).Create(article).Error
}

func (r *articleRepository) GetByID(ctx context.Context, id uint) (*entity.Article, error) {
	var article entity.Article
	err := r.db.WithContext(ctx).First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) GetAll(ctx context.Context) ([]*entity.Article, error) {
	var articles []*entity.Article
	err := r.db.WithContext(ctx).Find(&articles).Error
	return articles, err
}

func (r *articleRepository) GetPublished(ctx context.Context) ([]*entity.Article, error) {
	var articles []*entity.Article
	err := r.db.WithContext(ctx).Where("published = ?", true).Find(&articles).Error
	return articles, err
}

func (r *articleRepository) Update(ctx context.Context, article *entity.Article) error {
	return r.db.WithContext(ctx).Save(article).Error
}

func (r *articleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Article{}, id).Error
}

func (r *articleRepository) GetByAuthor(ctx context.Context, author string) ([]*entity.Article, error) {
	var articles []*entity.Article
	err := r.db.WithContext(ctx).Where("author = ?", author).Find(&articles).Error
	return articles, err
}