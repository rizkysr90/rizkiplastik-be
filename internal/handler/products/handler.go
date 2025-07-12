package products

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizkysr90/rizkiplastik-be/internal/util"
)

type Handler struct {
	service ProductService
}

func NewHandler(service ProductService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	endpoint := router.Group("/api/v1/products")
	endpoint.POST("/", h.CreateProduct)
	endpoint.PUT("/:product_id/single-product-type", h.UpdateSingleProductType)
	endpoint.PUT("/:product_id/variant-product-type", h.UpdateVariantProductType)
}

func (h *Handler) CreateProduct(c *gin.Context) {
	request := &CreateProductRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(c, request); err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) UpdateSingleProductType(c *gin.Context) {
	productID := c.Param("product_id")
	request := &UpdateSingleProductTypeRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.ProductID = productID
	if err := h.service.UpdateSingleProductType(c, request); err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
func (h *Handler) UpdateVariantProductType(c *gin.Context) {
	productID := c.Param("product_id")
	request := &UpdateVariantProductTypeRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.ProductID = productID
	if err := h.service.UpdateVariantProductType(c, request); err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
