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
	endpoint.GET("/", h.GetSizeUnits)
	endpoint.GET("/:size_unit_id", h.GetSizeUnit)
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
func (h *Handler) GetSizeUnits(c *gin.Context) {

	pageNumber := c.Query("page_number")
	pageSize := c.Query("page_size")
	sizeUnitName := c.Query("size_unit_name")
	sizeUnitCode := c.Query("size_unit_code")
	sizeUnitType := c.Query("size_unit_type")
	isActive := c.Query("is_active")

	pagination, err := util.NewPaginationData(pageNumber, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request := model.RequestGetSizeUnits{
		PaginationData: *pagination,
		SizeUnitName:   sizeUnitName,
		SizeUnitCode:   sizeUnitCode,
		SizeUnitType:   sizeUnitType,
		IsActive:       isActive,
	}

	response, err := h.service.GetSizeUnits(c, &request)
	if err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetSizeUnit(c *gin.Context) {
	sizeUnitID := c.Param("size_unit_id")
	request := model.RequestGetSizeUnit{
		SizeUnitID: sizeUnitID,
	}
	response, err := h.service.GetSizeUnit(c, &request)
	if err != nil {
		util.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}
