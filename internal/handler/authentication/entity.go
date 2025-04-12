package authentication

import (
	"time"
)

// Role represents available user roles
type Role string

const (
	// RoleAdmin represents admin role
	RoleAdmin Role = "ADMIN"
	// RoleStaff represents staff role
	RoleStaff Role = "STAFF"
)

// User represents a user entity
type User struct {
	Username     string     `json:"username"`
	HashPassword string     `json:"-"`
	Role         Role       `json:"role"`
	Token        *string    `json:"-"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
}
