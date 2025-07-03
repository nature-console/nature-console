package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nature-console/backend/internal/handler/auth"
	"github.com/nature-console/backend/internal/middleware"
	authUC "github.com/nature-console/backend/internal/usecase/auth"
)

func SetupAuthRoutes(router *gin.RouterGroup, handler *auth.Handler, authUseCase *authUC.UseCase) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", handler.Login)
		authRoutes.POST("/logout", handler.Logout)
		
		// Protected routes
		protected := authRoutes.Group("")
		protected.Use(middleware.AuthMiddleware(authUseCase))
		{
			protected.GET("/me", handler.Me)
		}
	}
}