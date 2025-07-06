package products

import "github.com/gin-gonic/gin"

type Handler struct {
	service ProductService
}

func NewHandler(service ProductService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// endpoint := router.Group("/api/v1/products")
}
