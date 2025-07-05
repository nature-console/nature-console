package auth

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/nature-console/backend/internal/domain/entity"
	authUC "github.com/nature-console/backend/internal/usecase/auth"
)

type Handler struct {
	authUseCase *authUC.UseCase
}

func NewHandler(authUseCase *authUC.UseCase) *Handler {
	return &Handler{
		authUseCase: authUseCase,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authUseCase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set JWT as HTTP-only cookie
	c.SetCookie(
		"token",
		response.Token,
		24*60*60, // 24 hours
		"/",
		"",
		false, // Set to true in production with HTTPS
		true,  // HTTP-only
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    response.User,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie(
		"token",
		"",
		-1, // Expire immediately
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

func (h *Handler) Me(c *gin.Context) {
	// Get user from middleware context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Type assertion to ensure user is of correct type
	user, ok := userInterface.(*entity.AdminUser)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}