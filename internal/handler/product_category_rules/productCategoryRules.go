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
	endpoint := router.Group("/api/v1/categories-rules")

	endpoint.POST("/:product_category_id/packaging-rules", h.PostPackagingRules)
	endpoint.PUT("/:product_category_id/packaging-rules/:rule_id", h.PutPackagingRules)
	endpoint.GET("/:product_category_id/packaging-rules", h.GetPackagingRules)
	endpoint.PATCH("/packaging-rules/:rule_id/status", h.PutPackagingRulesStatus)

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

func (h *Handler) PutPackagingRules(c *gin.Context) {
	productCategoryID := c.Param("product_category_id")
	ruleID := c.Param("rule_id")
	var request model.UpdateRulesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperror.NewBadRequest(c, httperror.WithMessage(err.Error()))
		return
	}
	request.ProductCategoryID = productCategoryID
	request.RuleID = ruleID
	if err := h.service.UpdateRules(c, &request); err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) GetPackagingRules(c *gin.Context) {
	productCategoryID := c.Param("product_category_id")
	status := c.Query("status")
	var request model.GetListRulesRequest
	request.ProductCategoryID = productCategoryID
	request.Status = status
	response, err := h.service.GetRules(c, &request)
	if err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}
func (h *Handler) PutPackagingRulesStatus(c *gin.Context) {
	ruleID := c.Param("rule_id")
	var request model.UpdateRulesStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperror.NewBadRequest(c, httperror.WithMessage(err.Error()))
		return
	}
	request.RuleID = ruleID
	if err := h.service.UpdateStatusRules(c, &request); err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
