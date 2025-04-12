package authentication

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/config"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *pgxpool.Pool, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

// RegisterRoutes registers all authentication-related routes
func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1/auth")
	{
		api.POST("/login", h.Login)
	}
}

// Login handles user authentication and token generation
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Step 1 & 2: Check existence of the username and get stored hash
	var user User
	var hashPassword string
	var role string
	query := `
		SELECT username, hash_password, role
		FROM users
		WHERE username = $1
	`

	err := h.db.QueryRow(c, query, req.Username).Scan(
		&user.Username,
		&hashPassword,
		&role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate user"})
		return
	}
	// Step 3: Compare password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	// Step 4: Create JWT token with 7-day expiration
	token, err := h.generateJWTToken(user.Username, Role(role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	// Update token in database for manual revocation capability
	updateQuery := `
		UPDATE users
		SET token = $1, last_login = $2
		WHERE username = $3
	`

	now := time.Now().UTC()
	_, err = h.db.Exec(c, updateQuery, token, now, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update token"})
		return
	}
	// Step 5: Return token to the response body
	response := LoginResponse{
		Token: token,
	}

	c.JSON(http.StatusOK, response)
}

// generateJWTToken creates a new JWT token with specified claims
func (h *AuthHandler) generateJWTToken(username string, role Role) (string, error) {
	// Get secret key from environment variable or use default for development
	secretKey := h.cfg.JWTSecret
	if secretKey == "" {
		return "", errors.New("jwt secret cannot be empty")
	}

	// Create token with claims
	claims := jwt.MapClaims{
		"username": username,
		"role":     string(role),
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days expiration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
