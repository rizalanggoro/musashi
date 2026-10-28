package auth

import (
	"github.com/google/uuid"
	"musashi/backend/internal/domain"
)

type SessionPayload struct {
	UserID    uuid.UUID       `json:"sub"`
	Email     string          `json:"email"`
	Name      string          `json:"name"`
	Role      domain.UserRole `json:"role"`
	ExpiresAt int64           `json:"exp"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    string          `json:"id"`
	Email string          `json:"email"`
	Name  string          `json:"name"`
	Role  domain.UserRole `json:"role"`
}
