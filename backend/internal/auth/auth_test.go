package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestAdminPasswordHash(t *testing.T) {
	// Bcrypt hash for "admin123"
	hash := "$2a$10$i/3zZPOBK/Bw79ayDXks9u3YblN8XnEAQqruvWzi2uwzl4u/TQerS"
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("admin123")); err != nil {
		t.Fatalf("Password hash does not match 'admin123': %v", err)
	}
}
