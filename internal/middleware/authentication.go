package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/config"
)

// AuthMiddleware provides JWT authentication middleware
type AuthMiddleware struct {
	db     *pgxpool.Pool
	config *config.Config
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(db *pgxpool.Pool, config *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		db:     db,
		config: config,
	}
}

// RequireAuth middleware checks for valid JWT token
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		// Extract token from "Bearer <token>"
		tokenString := ""
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString = parts[1]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}
		// Parse and validate token
		claims, err := m.validateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		tokenRole, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims: missing role"})
			c.Abort()
			return
		}
		// Get token and role from database by username
		var storedToken string
		var dbRole string
		query := `
			SELECT token, role
			FROM users
			WHERE username = $1
		`
		err = m.db.QueryRow(c, query, username).Scan(&storedToken, &dbRole)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found or token not set"})
			c.Abort()
			return
		}
		// Step 2: Compare stored token with provided token
		if storedToken != tokenString {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			c.Abort()
			return
		}
		// Verify that token role matches database role
		if tokenRole != dbRole {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token role does not match user role in database"})
			c.Abort()
			return
		}
		// Set claims in context for use in handlers
		c.Set("username", username)
		c.Set("role", tokenRole)
		c.Next()
	}
}

// RequireRole middleware checks if user has the required role
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context (set by RequireAuth middleware)
		roleValue, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found in context"})
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid role type"})
			c.Abort()
			return
		}

		// Check if user has one of the required roles
		hasValidRole := false
		for _, r := range roles {
			if role == r {
				hasValidRole = true
				break
			}
		}

		if !hasValidRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}
		c.Next()
	}
}
func (m *AuthMiddleware) validateToken(tokenString string) (jwt.MapClaims, error) {
	// Get secret key from environment variable or use default for development
	secretKey := m.config.JWTSecret
	if secretKey == "" {
		return nil, errors.New("empty secret")
	}
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Additional check to ensure we're using HS256 specifically
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: expected HS256, got %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	// Validate and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
