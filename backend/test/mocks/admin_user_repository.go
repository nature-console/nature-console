package mocks

import (
	"context"
	"errors"

	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/internal/domain/repository"
)

// MockAdminUserRepository is a mock implementation of AdminUserRepository for testing
type MockAdminUserRepository struct {
	adminUsers []*entity.AdminUser
	nextID     uint
}

// NewMockAdminUserRepository creates a new mock admin user repository
func NewMockAdminUserRepository() repository.AdminUserRepository {
	return &MockAdminUserRepository{
		adminUsers: make([]*entity.AdminUser, 0),
		nextID:     1,
	}
}

func (m *MockAdminUserRepository) Create(ctx context.Context, adminUser *entity.AdminUser) error {
	if adminUser == nil {
		return errors.New("admin user cannot be nil")
	}
	
	adminUser.ID = m.nextID
	m.nextID++
	
	m.adminUsers = append(m.adminUsers, adminUser)
	return nil
}

func (m *MockAdminUserRepository) GetByID(ctx context.Context, id uint) (*entity.AdminUser, error) {
	for _, user := range m.adminUsers {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("admin user not found")
}

func (m *MockAdminUserRepository) GetByEmail(ctx context.Context, email string) (*entity.AdminUser, error) {
	for _, user := range m.adminUsers {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("admin user not found")
}

func (m *MockAdminUserRepository) Update(ctx context.Context, adminUser *entity.AdminUser) error {
	if adminUser == nil {
		return errors.New("admin user cannot be nil")
	}
	
	for i, user := range m.adminUsers {
		if user.ID == adminUser.ID {
			m.adminUsers[i] = adminUser
			return nil
		}
	}
	return errors.New("admin user not found")
}

func (m *MockAdminUserRepository) Delete(ctx context.Context, id uint) error {
	for i, user := range m.adminUsers {
		if user.ID == id {
			m.adminUsers = append(m.adminUsers[:i], m.adminUsers[i+1:]...)
			return nil
		}
	}
	return errors.New("admin user not found")
}

func (m *MockAdminUserRepository) List(ctx context.Context) ([]*entity.AdminUser, error) {
	result := make([]*entity.AdminUser, len(m.adminUsers))
	copy(result, m.adminUsers)
	return result, nil
}