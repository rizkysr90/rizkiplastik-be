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
