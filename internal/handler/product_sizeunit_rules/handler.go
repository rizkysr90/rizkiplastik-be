package productsizeunitrules

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/product_sizeunit_rules/service"
	"github.com/rizkysr90/rizkiplastik-be/internal/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util"
)

type Handler struct {
	service *service.ProductSizeUnitRulesService
}

func NewHandler(productSizeUnitRulesRepository repository.ProductSizeUnitRules) *Handler {
	service := service.NewProductSizeUnitRulesService(productSizeUnitRulesRepository)
	return &Handler{service: service}
}
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	endpoint := router.Group("/api/v1/categories-rules")

	endpoint.POST("/:product_category_id/size-unit-rules", h.PostSizeUnitRules)
}

func (h *Handler) PostSizeUnitRules(c *gin.Context) {
	productCategoryID := c.Param("product_category_id")
	var request model.CreateSizeUnitRulesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	request.ProductCategoryID = productCategoryID
	if err := h.service.CreateRule(c, &request); err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}
