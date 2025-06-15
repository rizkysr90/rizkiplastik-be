package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

// HandleServiceError handles service errors and responds with appropriate HTTP status
func HandleServiceError(c *gin.Context, err error) {
	if httpError, ok := err.(*httperror.HTTPError); ok {
		switch httpError.Code {
		case http.StatusBadRequest:
			c.JSON(http.StatusBadRequest, httpError)
		case http.StatusNotFound:
			c.JSON(http.StatusNotFound, httpError)
		case http.StatusUnauthorized:
			c.JSON(http.StatusUnauthorized, httpError)
		case http.StatusForbidden:
			c.JSON(http.StatusForbidden, httpError)
		default:
			c.JSON(http.StatusInternalServerError, httpError)
		}
		return
	}
	if multiFieldValidation, ok := err.(*httperror.MultiFieldValidation); ok {
		c.JSON(http.StatusBadRequest, multiFieldValidation)
		return
	}

	// For non-ServiceError types, return 500
	c.JSON(http.StatusInternalServerError, httperror.HTTPError{
		Code:    http.StatusInternalServerError,
		Info:    "internal server error",
		Message: err.Error(),
	})
}
