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

	endpoint := router.Group("/api/v1/categories")
	{
		endpoint.POST("/", h.CreateCategory)
		endpoint.PUT("/:category_id", h.UpdateCategory)
		endpoint.GET("/", h.GetListCategory)
		endpoint.GET("/:category_id", h.GetByCategoryID)
	}
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var requestBody CreateCategoryRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.categoryService.CreateCategory(c, &requestBody)
	if err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}
func (h *Handler) UpdateCategory(c *gin.Context) {
	categoryID := c.Param("category_id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category_id is required"})
		return
	}
	var requestBody UpdateCategoryRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	requestBody.CategoryID = categoryID
	err := h.categoryService.UpdateCategory(c, &requestBody)
	if err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
func (h *Handler) GetListCategory(c *gin.Context) {
	pageSize := c.Query("page_size")
	pageNumber := c.Query("page_number")
	pagination, err := util.NewPaginationData(pageNumber, pageSize)
	if err != nil {
		errMsg := "invalid pagination data : " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}
	requestQueryParams := GetListCategoryRequest{
		PaginationData: *pagination,
		CategoryName:   c.Query("category_name"),
		CategoryCode:   c.Query("category_code"),
		IsActive:       c.Query("is_active"),
	}
	response, err := h.categoryService.GetListCategory(c, &requestQueryParams)
	if err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)

}
func (h *Handler) GetByCategoryID(c *gin.Context) {
	categoryID := c.Param("category_id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category_id is required"})
		return
	}
	response, err := h.categoryService.GetByCategoryID(c, &GetByCategoryIDRequest{CategoryID: categoryID})
	if err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}
