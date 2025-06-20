package packagingtypes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/packagingtypes/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/packagingtypes/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/packagingtypes/service"
	"github.com/rizkysr90/rizkiplastik-be/internal/util"
)

// PackagingTypesHandler handles HTTP requests for packaging types
type Handler struct {
	service *service.PackagingType
}

// NewHandler creates a new packaging types handler
func NewHandler(packagingTypeRepo repository.PackagingType) *Handler {
	service := service.NewPackagingType(packagingTypeRepo)
	return &Handler{service: &service}
}

// RegisterRoutes registers all category related routes
func (h *Handler) RegisterRoutes(
	router *gin.Engine) {

	endpoint := router.Group("/api/v1/packaging-types")
	{
		endpoint.POST("/", h.PostPackagingType)
		endpoint.PUT("/:packaging_type_id", h.UpdatePackagingType)
		endpoint.GET("/", h.GetPaginatedPackagingTypes)
		endpoint.GET("/:packaging_type_id", h.GetPackagingType)
	}
}
func (h *Handler) PostPackagingType(c *gin.Context) {
	var req model.RequestCreatePackagingType
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.CreatePackagingType(c, &req)
	if err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

func (h *Handler) UpdatePackagingType(c *gin.Context) {
	packagingTypeID := c.Param("packaging_type_id")
	var req model.RequestUpdatePackagingType
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.PackagingID = packagingTypeID
	err := h.service.UpdatePackagingType(c, &req)
	if err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) GetPaginatedPackagingTypes(c *gin.Context) {
	pageSize := c.Query("page_size")
	pageNumber := c.Query("page_number")
	pagination, err := util.NewPaginationData(pageNumber, pageSize)
	if err != nil {
		util.HandleServiceError(c, err)
		return
	}
	requestQueryParams := model.RequestGetPackagingTypes{
		PaginationData: *pagination,
		Name:           c.Query("name"),
		Code:           c.Query("code"),
		IsActive:       c.Query("is_active"),
	}
	response, err := h.service.GetPackagingTypes(c, &requestQueryParams)
	if err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}
func (h *Handler) GetPackagingType(c *gin.Context) {
	packagingTypeID := c.Param("packaging_type_id")
	request := model.RequestGetPackagingType{
		PackagingID: packagingTypeID,
	}
	response, err := h.service.GetPackagingType(c, &request)
	if err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}
