package category

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/middleware"
	"github.com/rizkysr90/rizkiplastik-be/internal/util"
)

// CategoryHandler handles HTTP requests for products
type Handler struct {
	db              *pgxpool.Pool
	categoryService Service
}

// CategoryHandler creates a new product handler
func NewCategoryHandler(db *pgxpool.Pool, categoryService Service) *Handler {
	return &Handler{db: db, categoryService: categoryService}
}

// RegisterRoutes registers all category related routes
func (h *Handler) RegisterRoutes(
	router *gin.Engine,
	authMiddleware *middleware.AuthMiddleware) {

	// Admin routes (only admin can access)
	admin := router.Group("/api/v1/categories")
	admin.Use(authMiddleware.RequireAuth(), authMiddleware.RequireRole("ADMIN"))
	{
		admin.POST("", h.CreateCategory)
	}
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var requestBody CreateCategoryRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.categoryService.CreateCategory(c.Request.Context(), &requestBody)
	if err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}
