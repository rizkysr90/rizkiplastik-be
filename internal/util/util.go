package util

import "github.com/gin-gonic/gin"

// GetUsernameFromContext extracts the username from the provided context.
// It returns the username as a string and a boolean indicating if the extraction was successful.
func GetUsernameFromContext(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}

	usernameStr, ok := username.(string)
	if !ok {
		return "", false
	}

	return usernameStr, true
}
