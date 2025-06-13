package util

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository/pg"
)

// ServiceError represents a service layer error with HTTP status code
type ServiceError struct {
	HTTPCode int
	Message  string
}

func (e *ServiceError) Error() string {
	return e.Message
}

// HandleServiceError handles service errors and responds with appropriate HTTP status
func HandleServiceError(c *gin.Context, err error) {
	if serviceErr, ok := err.(*ServiceError); ok {
		switch serviceErr.HTTPCode {
		case http.StatusBadRequest:
			c.JSON(http.StatusBadRequest, gin.H{"error": serviceErr.Message})
		case http.StatusNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": serviceErr.Message})
		case http.StatusConflict:
			c.JSON(http.StatusConflict, gin.H{"error": serviceErr.Message})
		case http.StatusUnauthorized:
			c.JSON(http.StatusUnauthorized, gin.H{"error": serviceErr.Message})
		case http.StatusForbidden:
			c.JSON(http.StatusForbidden, gin.H{"error": serviceErr.Message})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": serviceErr.Message})
		}
		return
	}

	// For non-ServiceError types, return 500
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// ConvertRepositoryError converts repository errors to ServiceError
func ConvertRepositoryError(err error) *ServiceError {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, pg.ErrCategoryAlreadyExists):
		return &ServiceError{
			HTTPCode: 400,
			Message:  "category already exists",
		}
	case errors.Is(err, pg.ErrCategoryNotFound):
		return &ServiceError{
			HTTPCode: 404,
			Message:  "category not found",
		}
	case errors.Is(err, pg.ErrDatabaseOperation):
		return &ServiceError{
			HTTPCode: 500,
			Message:  "internal server error" + err.Error(),
		}
	case errors.Is(err, pg.ErrTransactionFailed):
		return &ServiceError{
			HTTPCode: 500,
			Message:  "internal server error" + err.Error(),
		}
	default:
		// For any unknown repository error, return 500
		return &ServiceError{
			HTTPCode: 500,
			Message:  "internal server error" + err.Error(),
		}
	}
}
