package authentication

// LoginRequest represents data needed for login
type LoginRequest struct {
	Username string `json:"username" binding:"required,max=30"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response for login endpoint
type LoginResponse struct {
	Token string `json:"token"`
}

// UserClaims represents the JWT claims for user authentication
type UserClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
