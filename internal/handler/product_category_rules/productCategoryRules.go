package productcategoryrules

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/product_category_rules/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/product_category_rules/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/product_category_rules/service"
	"github.com/rizkysr90/rizkiplastik-be/internal/util"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type Handler struct {
	service *service.ProductCategoryRules
}

func NewHandler(productCategoryRules repository.ProductCategoryRules) *Handler {
	service := service.NewProductCategoryRules(productCategoryRules)
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	endpoint := router.Group("/api/v1/categories/:product_category_id/packaging-rules")

	endpoint.POST("/", h.PostPackagingRules)
}

func (h *Handler) PostPackagingRules(c *gin.Context) {
	productCategoryID := c.Param("product_category_id")
	var request model.CreateRulesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperror.NewBadRequest(c, httperror.WithMessage(err.Error()))
		return
	}
	request.ProductCategoryID = productCategoryID
	if err := h.service.CreateRules(c, &request); err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}
