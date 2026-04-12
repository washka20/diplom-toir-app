package services

import (
	"context"
	"fmt"
	"time"

	"toir-app/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// TokenPair содержит пару JWT токенов: access и refresh.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// AuthService инкапсулирует бизнес-логику аутентификации.
type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

// NewAuthService создаёт сервис аутентификации.
func NewAuthService(repo repository.UserRepository, secret string) *AuthService {
	return &AuthService{
		userRepo:  repo,
		jwtSecret: secret,
	}
}

// Login проверяет учётные данные и возвращает пару JWT токенов.
func (s *AuthService) Login(ctx context.Context, username, password string) (*TokenPair, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("user account is deactivated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	accessToken, err := s.generateToken(user.ID, user.Role, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateToken(user.ID, user.Role, 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) generateToken(userID uint, role string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": float64(userID),
		"role":    role,
		"exp":     jwt.NewNumericDate(time.Now().Add(ttl)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
