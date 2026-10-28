package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"musashi/backend/internal/domain"
)

const (
	CookieName   = "musashi_session"
	CookieMaxAge = 7 * 24 * 3600 // 7 days
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidToken       = errors.New("invalid token format")
	ErrInvalidSignature   = errors.New("invalid signature")
)

var sessionSecret = []byte("musashi-super-secret-session-key-change-in-prod")

type Service interface {
	Login(input LoginInput) (*domain.User, string, error)
	ParseToken(token string) (*SessionPayload, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func sign(data []byte) string {
	mac := hmac.New(sha256.New, sessionSecret)
	mac.Write(data)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func createToken(user *domain.User) (string, error) {
	payload := SessionPayload{
		UserID:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString(bytes)
	signature := sign([]byte(encoded))
	return encoded + "." + signature, nil
}

func parseToken(token string) (*SessionPayload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, ErrInvalidToken
	}
	expectedSig := sign([]byte(parts[0]))
	if !hmac.Equal([]byte(expectedSig), []byte(parts[1])) {
		return nil, ErrInvalidSignature
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	var payload SessionPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	if time.Now().Unix() > payload.ExpiresAt {
		return nil, ErrTokenExpired
	}
	return &payload, nil
}

func (s *service) Login(input LoginInput) (*domain.User, string, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := createToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *service) ParseToken(token string) (*SessionPayload, error) {
	return parseToken(token)
}
