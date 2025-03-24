package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type JWTTestSuite struct {
	suite.Suite
	service JWTTokenServicer
}

func (suite *JWTTestSuite) SetupTest() {
	cfg := viper.New()
	cfg.Set("auth.jwt_secret", "test_secret_key")
	cfg.Set("auth.jwt_expiration", time.Hour)
	suite.service = NewJWTTokenService(cfg)
}

func (suite *JWTTestSuite) TestGenerateToken() {
	userID := uuid.New()
	email := "test@example.com"
	role := "employee"

	token, err := suite.service.GenerateToken(userID, email, role)
	suite.NoError(err)
	suite.NotEmpty(token)

	// Validate the generated token
	claims, err := suite.service.ValidateToken(token)
	suite.NoError(err)
	suite.Equal(userID, claims.UserID)
	suite.Equal(email, claims.Email)
	suite.Equal(role, claims.Role)
}

func (suite *JWTTestSuite) TestValidateToken() {
	// Generate a valid token
	userID := uuid.New()
	email := "test@example.com"
	role := "employee"
	token, err := suite.service.GenerateToken(userID, email, role)
	suite.NoError(err)

	// Test valid token
	claims, err := suite.service.ValidateToken(token)
	suite.NoError(err)
	suite.Equal(userID, claims.UserID)
	suite.Equal(email, claims.Email)
	suite.Equal(role, claims.Role)

	// Test invalid token
	claims, err = suite.service.ValidateToken("invalid_token")
	suite.Error(err)
	suite.Equal(ErrInvalidToken, err)
}

func TestJWTTestSuite(t *testing.T) {
	suite.Run(t, new(JWTTestSuite))
}
