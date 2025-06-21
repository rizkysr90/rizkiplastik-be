package middleware

import (
	"github.com/gin-gonic/gin"
)

func Session() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", "rizkysr90")
		c.Set("role", "admin")
		c.Next()
	}
}
