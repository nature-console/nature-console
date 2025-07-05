package admin

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	articleUC "github.com/nature-console/backend/internal/usecase/article"
)

type Handler struct {
	articleUseCase *articleUC.UseCase
}

func NewHandler(articleUseCase *articleUC.UseCase) *Handler {
	return &Handler{
		articleUseCase: articleUseCase,
	}
}

type DashboardStats struct {
	TotalArticles     int64 `json:"total_articles"`
	PublishedArticles int64 `json:"published_articles"`
	DraftArticles     int64 `json:"draft_articles"`
}

func (h *Handler) GetDashboard(c *gin.Context) {
	// Get all articles to calculate stats
	allArticles, err := h.articleUseCase.GetAllArticles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	publishedArticles, err := h.articleUseCase.GetPublishedArticles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stats := DashboardStats{
		TotalArticles:     int64(len(allArticles)),
		PublishedArticles: int64(len(publishedArticles)),
		DraftArticles:     int64(len(allArticles)) - int64(len(publishedArticles)),
	}

	// Get recent articles (last 5)
	recentArticles := allArticles
	if len(allArticles) > 5 {
		recentArticles = allArticles[:5]
	}

	c.JSON(http.StatusOK, gin.H{
		"stats":           stats,
		"recent_articles": recentArticles,
	})
}

func (h *Handler) GetAllArticles(c *gin.Context) {
	articles, err := h.articleUseCase.GetAllArticles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, articles)
}

type CreateArticleRequest struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Author    string `json:"author" binding:"required"`
	Published bool   `json:"published"`
}

func (h *Handler) CreateArticle(c *gin.Context) {
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article, err := h.articleUseCase.CreateArticle(c.Request.Context(), req.Title, req.Content, req.Author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If published is true, publish the article
	if req.Published {
		article, err = h.articleUseCase.PublishArticle(c.Request.Context(), article.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, article)
}

type UpdateArticleRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	Published bool   `json:"published"`
}

func (h *Handler) UpdateArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}

	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article, err := h.articleUseCase.UpdateArticle(c.Request.Context(), uint(id), req.Title, req.Content, req.Author, req.Published)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}

func (h *Handler) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}

	err = h.articleUseCase.DeleteArticle(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) GetArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}

	article, err := h.articleUseCase.GetArticle(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}