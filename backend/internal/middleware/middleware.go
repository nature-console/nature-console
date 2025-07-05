package middleware

import (
	"github.com/gin-gonic/gin"
	authUC "github.com/nature-console/backend/internal/usecase/auth"
)

// SetupMiddlewares applies all common middlewares to the router
func SetupMiddlewares(r *gin.Engine) {
	r.Use(CORSMiddleware())
	r.Use(LoggingMiddleware())
}

// NewAuthMiddleware creates an auth middleware with the given use case
func NewAuthMiddleware(authUseCase *authUC.UseCase) gin.HandlerFunc {
	return AuthMiddleware(authUseCase)
}