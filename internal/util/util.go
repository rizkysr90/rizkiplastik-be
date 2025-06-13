package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

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

// ConvertStringToInt converts a string to an integer.
// It returns the integer and an error if the conversion fails.
func ConvertStringToInt(str string) (int, error) {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return num, nil
}
