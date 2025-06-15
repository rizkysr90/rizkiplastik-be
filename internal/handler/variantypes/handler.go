package variantypes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/variantypes/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/variantypes/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/variantypes/service"
	"github.com/rizkysr90/rizkiplastik-be/internal/util"
)

type Handler struct {
	service *service.VarianTypes
}

func NewHandler(
	repository repository.VarianType,
) *Handler {
	service := service.NewVarianTypes(repository)
	return &Handler{
		service: service,
	}
}
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	endpoint := router.Group("/api/v1/variant-types")
	endpoint.POST("/", h.PostVariantType)
	endpoint.PUT("/:variant_type_id", h.PutVariantType)
	endpoint.GET("/", h.GetVarianTypes)
	endpoint.GET("/:variant_type_id", h.GetVarianType)
}

func (h *Handler) PostVariantType(c *gin.Context) {
	var input model.RequestCreateVarianType
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateVarianType(c, &input); err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}
func (h *Handler) PutVariantType(c *gin.Context) {
	variantTypeID := c.Param("variant_type_id")
	var input model.RequestUpdateVarianType
	input.VarianTypeID = variantTypeID
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateVarianType(c, &input); err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
func (h *Handler) GetVarianTypes(c *gin.Context) {
	pageNumber := c.Query("page_number")
	pageSize := c.Query("page_size")
	variantTypeName := c.Query("variant_type_name")
	isActive := c.Query("is_active")

	pagination, err := util.NewPaginationData(pageNumber, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request := model.RequestVarianTypePaginated{
		PaginationData: *pagination,
		VarianTypeName: variantTypeName,
		IsActive:       isActive,
	}
	response, err := h.service.GetVarianTypes(c, &request)
	if err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetVarianType(c *gin.Context) {
	variantTypeID := c.Param("variant_type_id")
	request := model.RequestGetVarianType{
		VarianTypeID: variantTypeID,
	}
	response, err := h.service.GetVarianType(c, &request)
	if err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}
