package auth

import (
	"musashi/backend/internal/domain"
)

type SessionPayload struct {
	UserID    uint            `json:"sub"`
	Email     string          `json:"email"`
	Name      string          `json:"name"`
	Role      domain.UserRole `json:"role"`
	ExpiresAt int64           `json:"exp"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID    string          `json:"id"`
	Email string          `json:"email"`
	Name  string          `json:"name"`
	Role  domain.UserRole `json:"role"`
	Token string          `json:"token,omitempty"`
}

type UserResponse = User
