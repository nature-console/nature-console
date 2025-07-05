package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nature-console/backend/internal/handler/admin"
	"github.com/nature-console/backend/internal/middleware"
	authUC "github.com/nature-console/backend/internal/usecase/auth"
)

func SetupAdminRoutes(router *gin.RouterGroup, handler *admin.Handler, authUseCase *authUC.UseCase) {
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(authUseCase))
	{
		adminRoutes.GET("/dashboard", handler.GetDashboard)
		
		// Admin article management
		articles := adminRoutes.Group("/articles")
		{
			articles.GET("", handler.GetAllArticles)
			articles.GET("/:id", handler.GetArticle)
			articles.POST("", handler.CreateArticle)
			articles.PUT("/:id", handler.UpdateArticle)
			articles.DELETE("/:id", handler.DeleteArticle)
		}
	}
}