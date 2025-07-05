package admin_user

import (
	"context"
	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type adminUserRepository struct {
	db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) repository.AdminUserRepository {
	return &adminUserRepository{db: db}
}

func (r *adminUserRepository) GetByEmail(ctx context.Context, email string) (*entity.AdminUser, error) {
	var user entity.AdminUser
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *adminUserRepository) GetByID(ctx context.Context, id uint) (*entity.AdminUser, error) {
	var user entity.AdminUser
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *adminUserRepository) Create(ctx context.Context, user *entity.AdminUser) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *adminUserRepository) Update(ctx context.Context, user *entity.AdminUser) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *adminUserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.AdminUser{}, id).Error
}

func (r *adminUserRepository) List(ctx context.Context) ([]*entity.AdminUser, error) {
	var users []*entity.AdminUser
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}