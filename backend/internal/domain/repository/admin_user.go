package repository

import (
	"context"
	"github.com/nature-console/backend/internal/domain/entity"
)

type AdminUserRepository interface {
	GetByEmail(ctx context.Context, email string) (*entity.AdminUser, error)
	GetByID(ctx context.Context, id uint) (*entity.AdminUser, error)
	Create(ctx context.Context, user *entity.AdminUser) error
	Update(ctx context.Context, user *entity.AdminUser) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]*entity.AdminUser, error)
}