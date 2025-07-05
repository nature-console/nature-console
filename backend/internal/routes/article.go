package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nature-console/backend/internal/handler/article"
)

func SetupArticleRoutes(router *gin.RouterGroup, handler *article.Handler) {
	articles := router.Group("/articles")
	{
		// Public routes - only published articles
		articles.GET("", handler.GetPublishedArticles) // Show only published articles for public
		articles.GET("/:id", handler.GetArticle)       // Get specific article (will need to check if published)
	}
}