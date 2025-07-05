package auth

import (
	"context"
	"errors"
	"github.com/nature-console/backend/internal/domain/entity"
	"github.com/nature-console/backend/internal/domain/repository"
	"github.com/nature-console/backend/internal/utils"
)

type UseCase struct {
	adminUserRepo repository.AdminUserRepository
}

func NewUseCase(adminUserRepo repository.AdminUserRepository) *UseCase {
	return &UseCase{
		adminUserRepo: adminUserRepo,
	}
}

type LoginResponse struct {
	Token string           `json:"token"`
	User  *entity.AdminUser `json:"user"`
}

func (u *UseCase) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := u.adminUserRepo.GetByEmail(ctx, email)
	if err != nil {
		// Return generic "invalid credentials" for any error to avoid user enumeration
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// Create user response without password hash
	userResponse := &entity.AdminUser{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &LoginResponse{
		Token: token,
		User:  userResponse,
	}, nil
}

func (u *UseCase) GetUserFromToken(ctx context.Context, tokenString string) (*entity.AdminUser, error) {
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	user, err := u.adminUserRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}