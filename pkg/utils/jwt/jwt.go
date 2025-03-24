package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// Common JWT errors
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// JWTTokenServicer defines the interface for JWT token operations
type JWTTokenServicer interface {
	GenerateToken(userID uuid.UUID, email string, role string) (string, error)
	ValidateToken(tokenString string) (*UserClaims, error)
}

// JWTTokenService handles JWT token generation and validation
type JWTTokenService struct {
	secretKey     []byte
	tokenDuration time.Duration
}

// NewJWTTokenService creates a new instance of JWTTokenService
func NewJWTTokenService(cfg *viper.Viper) JWTTokenServicer {
	return &JWTTokenService{
		secretKey:     []byte(cfg.GetString("auth.jwt_secret")),
		tokenDuration: cfg.GetDuration("auth.jwt_expiration"),
	}
}

// UserClaims represents the JWT claims for a user
type UserClaims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
}

// GenerateToken generates a new JWT token for a user
func (s *JWTTokenService) GenerateToken(userID uuid.UUID, email string, role string) (string, error) {
	// Create the claims
	now := time.Now()
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   userID.String(),
		},
		UserID: userID,
		Email:  email,
		Role:   role,
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	return token.SignedString(s.secretKey)
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTTokenService) ValidateToken(tokenString string) (*UserClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// Validate claims
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// Get and return the claims
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
