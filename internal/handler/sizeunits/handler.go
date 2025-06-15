package sizeunits

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/sizeunits/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/sizeunits/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/sizeunits/service"
	"github.com/rizkysr90/rizkiplastik-be/internal/util"
)

type Handler struct {
	service *service.SizeUnits
}

func NewHandler(repository repository.SizeUnits) *Handler {
	service := service.NewSizeUnits(repository)
	return &Handler{
		service: service,
	}
}
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	endpoint := router.Group("/api/v1/size-units")

	endpoint.POST("/", h.PostSizeUnit)
	endpoint.PUT("/:size_unit_id", h.PutSizeUnit)
}

func (h *Handler) PostSizeUnit(c *gin.Context) {
	var request model.RequestCreateSizeUnit
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateSizeUnit(c, request); err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

func (h *Handler) PutSizeUnit(c *gin.Context) {
	var request model.RequestUpdateSizeUnit
	request.SizeUnitID = c.Param("size_unit_id")
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateSizeUnit(c, request); err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
